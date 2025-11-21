import { cn } from "@/lib/utils"

interface CoverProps {
  src: string
  alt?: string
  size?: number
  className?: string
}

function Cover({ 
  src, 
  alt = "Cover", 
  size = 80, 
  className,
  ...props 
}: CoverProps) {
  // Si une className est fournie, on l'utilise sans style inline
  const shouldUseSize = !className?.includes('w-full') && !className?.includes('aspect-');
  const customStyle = shouldUseSize ? { width: `${size}px`, height: `${size}px` } : undefined;

  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-md flex items-center justify-center",
        shouldUseSize ? "" : "w-full",
        className
      )}
      style={customStyle}
      {...props}
    >
      <img
        src={src}
        alt={alt}
        className={cn(
          "object-cover",
          shouldUseSize ? "w-full h-full aspect-square" : "w-full h-full"
        )}
        loading="lazy"
      />
    </div>
  )
}

export { Cover, type CoverProps }