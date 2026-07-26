// Performance benchmarks for Rust Data Store

use criterion::{black_box, criterion_group, criterion_main, Criterion, BenchmarkId};
use std::time::Duration;

fn bench_set_operation(c: &mut Criterion) {
    let mut group = c.benchmark_group("set_operation");
    
    for size in [100, 1000, 10000].iter() {
        group.bench_with_input(BenchmarkId::from_parameter(size), size, |b, &size| {
            let store = DataStore::new();
            b.iter(|| {
                let key = format!("bench_key_{}", black_box(0));
                let value = vec![0u8; size];
                black_box(store.set(key.as_bytes(), value));
            });
        });
    }
    
    group.finish();
}

fn bench_get_operation(c: &mut Criterion) {
    let mut group = c.benchmark_group("get_operation");
    
    for size in [100, 1000, 10000].iter() {
        group.bench_with_input(BenchmarkId::from_parameter(size), size, |b, &size| {
            let store = DataStore::new();
            // Pre-populate
            for i in 0..1000 {
                let key = format!("bench_key_{}", i);
                let value = vec![0u8; size];
                store.set(key.as_bytes(), value);
            }
            
            b.iter(|| {
                let key = format!("bench_key_{}", black_box(0));
                black_box(store.get(key.as_bytes()));
            });
        });
    }
    
    group.finish();
}

fn bench_concurrent_operations(c: &mut Criterion) {
    let mut group = c.benchmark_group("concurrent_operations");
    
    group.bench_function("concurrent_sets", |b| {
        let store = DataStore::new();
        b.iter(|| {
            let key = format!("concurrent_key_{}", black_box(0));
            let value = vec![0u8; 100];
            black_box(store.set(key.as_bytes(), value));
        });
    });
    
    group.finish();
}

fn bench_ttl_eviction(c: &mut Criterion) {
    let mut group = c.benchmark_group("ttl_eviction");
    
    group.bench_function("ttl_cleanup", |b| {
        let store = DataStore::new();
        // Insert entries with TTL
        for i in 0..1000 {
            let key = format!("ttl_key_{}", i);
            let value = vec![0u8; 100];
            store.set_with_ttl(key.as_bytes(), value, Duration::from_secs(1));
        }
        
        b.iter(|| {
            black_box(store.cleanup_expired());
        });
    });
    
    group.finish();
}

fn bench_memory_usage(c: &mut Criterion) {
    let mut group = c.benchmark_group("memory_usage");
    
    group.bench_function("allocation", |b| {
        b.iter(|| {
            let data = vec![0u8; 1024];
            black_box(data);
        });
    });
    
    group.finish();
}

criterion_group! {
    name = data_store_benchmarks;
    config = Criterion::default()
        .measurement_time(Duration::from_secs(10))
        .sample_size(100);
    targets = bench_set_operation, bench_get_operation, bench_concurrent_operations, bench_ttl_eviction, bench_memory_usage
}

criterion_main!(data_store_benchmarks);
