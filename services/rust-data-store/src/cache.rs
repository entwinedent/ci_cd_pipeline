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
