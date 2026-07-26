use axum::{
    extract::{Path, State},
    http::StatusCode,
    response::Json,
    routing::{delete, get, put},
    Router,
};
use log::info;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tower_http::cors::CorsLayer;

mod cache;

#[derive(Clone)]
struct AppState {
    cache: Arc<cache::InMemoryCache>,
}

#[derive(Deserialize)]
struct SetValue {
    value: String,
    ttl_seconds: Option<i64>,
}

#[derive(Serialize)]
struct Response {
    success: bool,
    message: String,
    value: Option<String>,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    env_logger::Builder::from_default_env()
        .filter_level(log::LevelFilter::Info)
        .init();
    
    info!("Starting Rust Data Store service");
    
    // Initialize the in-memory cache
    let cache = Arc::new(cache::InMemoryCache::new());
    
    // Start background TTL cleanup task
    let cache_clone = Arc::clone(&cache);
    tokio::spawn(async move {
        cache_clone.cleanup_expired_entries().await;
    });
    
    let app_state = AppState { cache };
    
    let app = Router::new()
        .route("/healthz", get(health_check))
        .route("/api/v1/data/:key", get(get_data).put(set_data).delete(delete_data))
        .layer(CorsLayer::permissive())
        .with_state(app_state);
    
    let listener = tokio::net::TcpListener::bind("0.0.0.0:50051").await?;
    info!("Data Store HTTP server listening on 0.0.0.0:50051");
    
    axum::serve(listener, app).await?;
    
    Ok(())
}

async fn health_check() -> Json<Response> {
    Json(Response {
        success: true,
        message: "Service is healthy".to_string(),
        value: None,
    })
}

async fn get_data(
    State(state): State<AppState>,
    Path(key): Path<String>,
) -> Result<Json<Response>, StatusCode> {
    match state.cache.get(&key) {
        Some(value) => Ok(Json(Response {
            success: true,
            message: "Data retrieved successfully".to_string(),
            value: Some(String::from_utf8_lossy(&value).to_string()),
        })),
        None => Ok(Json(Response {
            success: false,
            message: "Key not found".to_string(),
            value: None,
        })),
    }
}

async fn set_data(
    State(state): State<AppState>,
    Path(key): Path<String>,
    Json(payload): Json<SetValue>,
) -> Result<Json<Response>, StatusCode> {
    state.cache.set(key, payload.value.into_bytes(), payload.ttl_seconds);
    Ok(Json(Response {
        success: true,
        message: "Data stored successfully".to_string(),
        value: None,
    }))
}

async fn delete_data(
    State(state): State<AppState>,
    Path(key): Path<String>,
) -> Result<Json<Response>, StatusCode> {
    let deleted = state.cache.delete(&key);
    Ok(Json(Response {
        success: deleted,
        message: if deleted {
            "Data deleted successfully".to_string()
        } else {
            "Key not found".to_string()
        },
        value: None,
    }))
}
