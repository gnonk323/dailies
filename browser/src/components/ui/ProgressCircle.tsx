const ProgressCircle = ({ value = 0, size = 100, strokeWidth = 2, colorClass = 'stroke-blue-500' }) => {
  // Ensure value stays between 0 and 1
  const clampedValue = Math.min(Math.max(value, 0), 1);
  
  // Calculate radius and circumference
  const center = size / 2;
  const radius = center - strokeWidth;
  const circumference = 2 * Math.PI * radius;
  
  // Calculate how much of the stroke to hide
  const strokeDashoffset = circumference - clampedValue * circumference;

  return (
    <svg 
      width={size} 
      height={size} 
      viewBox={`0 0 ${size} ${size}`}
      className="-rotate-90" // Rotates the whole SVG so it starts at 12 o'clock
    >
      {/* Background Track Circle */}
      <circle
        cx={center}
        cy={center}
        r={radius}
        fill="transparent"
        className="stroke-gray-200 dark:stroke-gray-700"
        strokeWidth={strokeWidth}
      />
      {/* Animated Progress Circle */}
      <circle
        cx={center}
        cy={center}
        r={radius}
        fill="transparent"
        className={`stroke-linecap-round transition-all duration-300 ease-in-out ${colorClass}`}
        strokeWidth={strokeWidth}
        strokeDasharray={circumference}
        strokeDashoffset={strokeDashoffset}
      />
    </svg>
  );
};

export default ProgressCircle;
