import React from 'react';

export default function Status() {
  return (
    <span className="pulse">
      <style jsx>{`
        .pulse {
          display: inline-block;
          width: 12px;
          height: 12px;
          border-radius: 50%;
          background: #cca92c;
          box-shadow: 0 0 0 rgba(204,169,44, 0.4);
          animation: pulse 2s infinite;
        }

        .pulse:hover {
          animation: none;
        }

        @keyframes pulse {
          0% {
            box-shadow: 0 0 0 0 rgba(204,169,44, 0.4);
          }
          70% {
            box-shadow: 0 0 0 6px rgba(204,169,44, 0);
          }
          100% {
            box-shadow: 0 0 0 0 rgba(204,169,44, 0);
          }
        }
      `}</style>
    </span>
  );
}
