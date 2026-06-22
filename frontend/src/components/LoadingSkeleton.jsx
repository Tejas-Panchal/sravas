// LoadingSkeleton renders animated placeholder cards while content loads
export default function LoadingSkeleton({ count = 8 }) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="animate-pulse">
          <div className="aspect-video bg-gray-800 rounded-lg mb-2" />
          <div className="flex gap-3">
            <div className="w-9 h-9 bg-gray-800 rounded-full shrink-0" />
            <div className="flex-1 space-y-2">
              <div className="h-3 bg-gray-800 rounded w-3/4" />
              <div className="h-3 bg-gray-800 rounded w-1/2" />
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
