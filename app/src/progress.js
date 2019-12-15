import React, { useMemo } from 'react';

export default function ({
  value,
}) {
  const radius = 12;
  const stroke = 3;

  const normalizedRadius = useMemo(() => {
    return radius - stroke / 2;
  }, [radius, stroke]);

  const style = useMemo(() => {
    const circumference = normalizedRadius * 2 * Math.PI;

    return {
      // strokeWidth: stroke,
      strokeDashoffset: circumference - value / 100 * circumference,
      strokeDasharray: circumference + ' ' + circumference,
      stroke: '#5335B8',
      fill: 'transparent',
    };
  }, [value, normalizedRadius]);

  return (
    <span className="spinner">
      <svg
        height={radius * 2}
        width={radius * 2}
      >
        <circle
          stroke="#f2f2f2"
          fill="transparent"
          strokeWidth={stroke}
          r={normalizedRadius}
          cx={radius}
          cy={radius}
        />
        <circle
          style={style}
          strokeWidth={stroke}
          r={normalizedRadius}
          cx={radius}
          cy={radius}
        />
      </svg>
      <style jsx>{`
        circle {
          transition: stroke-dashoffset 0.35s;
          transform: rotate(-90deg);
          transform-origin: 50% 50%;
        }
      `}</style>
    </span>
  );
}
