import { useEffect, useState } from "react"

interface SoundwaveProps {
  isPlaying?: boolean
  barCount?: number
  barColor?: string
  height?: number
}

export function Soundwave({ isPlaying = true, barCount = 5, barColor = "bg-primary", height = 24 }: SoundwaveProps) {
  const [barHeights, setBarHeights] = useState<number[]>(Array(barCount).fill(0.3))

  useEffect(() => {
    if (!isPlaying) {
      setBarHeights(Array(barCount).fill(0.3))
      return
    }

    const interval = setInterval(() => {
      setBarHeights((prev) =>
        prev.map(() => {
          return Math.random() * 0.8 + 0.2
        }),
      )
    }, 250)

    return () => clearInterval(interval)
  }, [isPlaying, barCount])

  return (
    <div className="flex items-center justify-center gap-1">
      {barHeights.map((heightPercent, index) => (
        <div
          key={index}
          className={`${barColor} rounded-full transition-all duration-150 ease-out`}
          style={{
            width: "4px",
            height: `${height * heightPercent}px`,
            opacity: 0.7,
          }}
          aria-hidden="true"
        />
      ))}
    </div>
  )
}
