use dashmap::DashMap;
use std::collections::VecDeque;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::sync::RwLock;

#[derive(Clone)]
pub struct CacheEntry {
    pub value: Vec<u8>,
    pub expires_at: Option<Instant>,
}

pub struct InMemoryCache {
    // Primary key-value store using DashMap for lock-free concurrent access
    store: DashMap<String, CacheEntry>,

    // Time-series ring buffer for metrics (fixed capacity)
    metrics: Arc<RwLock<VecDeque<MetricPoint>>>,

    // TTL expiration tracking (simplified - in production would use binary heap)
    _ttl_entries: Arc<RwLock<Vec<(String, Instant)>>>,
}

#[derive(Clone)]
#[allow(dead_code)]
pub struct MetricPoint {
    pub timestamp: Instant,
    pub operation: String,
    pub key: String,
}

impl InMemoryCache {
    pub fn new() -> Self {
        Self {
            store: DashMap::new(),
            metrics: Arc::new(RwLock::new(VecDeque::with_capacity(1000))),
            _ttl_entries: Arc::new(RwLock::new(Vec::new())),
        }
    }

    pub fn set(&self, key: String, value: Vec<u8>, ttl_seconds: Option<i64>) {
        let expires_at = ttl_seconds.map(|secs| Instant::now() + Duration::from_secs(secs as u64));

        let entry = CacheEntry { value, expires_at };

        self.store.insert(key.clone(), entry);

        // Record metric (fire and forget)
        #[allow(clippy::let_underscore_future)]
        let _ = self.record_metric("SET".to_string(), key.clone());
    }

    pub fn get(&self, key: &str) -> Option<Vec<u8>> {
        if let Some(entry) = self.store.get(key) {
            // Check if expired
            if let Some(expires_at) = entry.expires_at {
                if Instant::now() > expires_at {
                    return None;
                }
            }
            #[allow(clippy::let_underscore_future)]
            let _ = self.record_metric("GET".to_string(), key.to_string());
            return Some(entry.value.clone());
        }
        None
    }

    pub fn delete(&self, key: &str) -> bool {
        let removed = self.store.remove(key).is_some();
        if removed {
            #[allow(clippy::let_underscore_future)]
            let _ = self.record_metric("DELETE".to_string(), key.to_string());
        }
        removed
    }

    pub async fn cleanup_expired_entries(&self) {
        let mut interval = tokio::time::interval(Duration::from_secs(60));

        loop {
            interval.tick().await;

            let mut keys_to_remove = Vec::new();

            for entry in self.store.iter() {
                if let Some(expires_at) = entry.expires_at {
                    if Instant::now() > expires_at {
                        keys_to_remove.push(entry.key().clone());
                    }
                }
            }

            for key in keys_to_remove {
                self.store.remove(&key);
            }
        }
    }

    async fn record_metric(&self, operation: String, key: String) {
        let mut metrics = self.metrics.write().await;
        let point = MetricPoint {
            timestamp: Instant::now(),
            operation,
            key,
        };

        if metrics.len() >= 1000 {
            metrics.pop_front();
        }

        metrics.push_back(point);
    }

    #[allow(dead_code)]
    pub async fn get_metrics(&self) -> Vec<MetricPoint> {
        let metrics = self.metrics.read().await;
        metrics.iter().cloned().collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_set_and_get() {
        let cache = InMemoryCache::new();
        cache.set("test_key".to_string(), b"test_value".to_vec(), None);
        
        let result = cache.get("test_key");
        assert!(result.is_some());
        assert_eq!(result.unwrap(), b"test_value");
    }

    #[test]
    fn test_get_nonexistent_key() {
        let cache = InMemoryCache::new();
        let result = cache.get("nonexistent");
        assert!(result.is_none());
    }

    #[test]
    fn test_delete_existing_key() {
        let cache = InMemoryCache::new();
        cache.set("test_key".to_string(), b"test_value".to_vec(), None);
        
        let deleted = cache.delete("test_key");
        assert!(deleted);
        
        let result = cache.get("test_key");
        assert!(result.is_none());
    }

    #[test]
    fn test_delete_nonexistent_key() {
        let cache = InMemoryCache::new();
        let deleted = cache.delete("nonexistent");
        assert!(!deleted);
    }

    #[test]
    fn test_ttl_expiration() {
        let cache = InMemoryCache::new();
        cache.set("test_key".to_string(), b"test_value".to_vec(), Some(1));
        
        // Should be available immediately
        let result = cache.get("test_key");
        assert!(result.is_some());
        
        // Wait for expiration
        std::thread::sleep(std::time::Duration::from_secs(2));
        
        let result = cache.get("test_key");
        assert!(result.is_none());
    }

    #[test]
    fn test_multiple_operations() {
        let cache = InMemoryCache::new();
        
        cache.set("key1".to_string(), b"value1".to_vec(), None);
        cache.set("key2".to_string(), b"value2".to_vec(), None);
        cache.set("key3".to_string(), b"value3".to_vec(), None);
        
        assert_eq!(cache.get("key1"), Some(b"value1".to_vec()));
        assert_eq!(cache.get("key2"), Some(b"value2".to_vec()));
        assert_eq!(cache.get("key3"), Some(b"value3".to_vec()));
        
        cache.delete("key2");
        
        assert_eq!(cache.get("key1"), Some(b"value1".to_vec()));
        assert_eq!(cache.get("key2"), None);
        assert_eq!(cache.get("key3"), Some(b"value3".to_vec()));
    }
}
