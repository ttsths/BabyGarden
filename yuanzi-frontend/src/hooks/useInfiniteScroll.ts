import { useState, useEffect, useCallback, useRef } from 'react';

interface UseInfiniteScrollOptions {
  onLoadMore: () => Promise<void>;
  hasMore: boolean;
  threshold?: number;
}

/**
 * 无限滚动 Hook
 * @param options 配置选项
 */
export function useInfiniteScroll({ onLoadMore, hasMore, threshold = 100 }: UseInfiniteScrollOptions) {
  const [isLoading, setIsLoading] = useState(false);
  const observerRef = useRef<IntersectionObserver | null>(null);
  const loadMoreRef = useRef<HTMLDivElement | null>(null);

  const handleLoadMore = useCallback(async () => {
    if (isLoading || !hasMore) return;

    setIsLoading(true);
    try {
      await onLoadMore();
    } catch (error) {
      console.error('Load more failed:', error);
    } finally {
      setIsLoading(false);
    }
  }, [isLoading, hasMore, onLoadMore]);

  useEffect(() => {
    if (!hasMore) return;

    observerRef.current = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          handleLoadMore();
        }
      },
      {
        rootMargin: `${threshold}px`,
      }
    );

    if (loadMoreRef.current) {
      observerRef.current.observe(loadMoreRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [hasMore, handleLoadMore, threshold]);

  return {
    isLoading,
    loadMoreRef,
  };
}
