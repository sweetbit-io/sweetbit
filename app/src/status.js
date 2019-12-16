import React from 'react';

export default function Status({
  status,
}) {
  return (
    <span className={`pulse ${status}`}>
      <style jsx>{`
        .pulse {
          display: inline-block;
          width: 12px;
          height: 12px;
          border-radius: 50%;
          background: #c0c0c0;
          box-shadow: 0 0 0 rgba(192,192,192,0.4);
          animation: pulse-gray 2s infinite;
        }

        .pulse.green {
          background: #3cb371;
          box-shadow: 0 0 0 rgba(60,179,113,0.4);
          animation: pulse-green 2s infinite;
        }

        .pulse.amber {
          background: #cca92c;
          box-shadow: 0 0 0 rgba(204,169,44,0.4);
          animation: pulse-amber 2s infinite;
        }

        .pulse.red {
          background: #ff4500;
          box-shadow: 0 0 0 rgba(255,69,0,0.4);
          animation: pulse-red 2s infinite;
        }

        .pulse:hover {
          animation: none;
        }

        @keyframes pulse-gray {
          0% { box-shadow: 0 0 0 0 rgba(192,192,192,0.4); }
          70% { box-shadow: 0 0 0 6px rgba(192,192,192,0); }
          100% { box-shadow: 0 0 0 0 rgba(192,192,192,0); }
        }

        @keyframes pulse-green {
          0% { box-shadow: 0 0 0 0 rgba(60,179,113,0.4); }
          70% { box-shadow: 0 0 0 6px rgba(60,179,113,0); }
          100% { box-shadow: 0 0 0 0 rgba(60,179,113,0); }
        }

        @keyframes pulse-amber {
          0% { box-shadow: 0 0 0 0 rgba(204,169,44,0.4); }
          70% { box-shadow: 0 0 0 6px rgba(204,169,44,0); }
          100% { box-shadow: 0 0 0 0 rgba(204,169,44,0); }
        }

        @keyframes pulse-red {
          0% { box-shadow: 0 0 0 0 rgba(255,69,0,0.4); }
          70% { box-shadow: 0 0 0 6px rgba(255,69,0,0); }
          100% { box-shadow: 0 0 0 0 rgba(255,69,0,0); }
        }
      `}</style>
    </span>
  );
}
