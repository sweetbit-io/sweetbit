import React from 'react';

const twelve = [...Array(12).keys()];

export default function () {
  return (
    <span className="spinner">
      {twelve.map((i) => <span key={i} />)}
      <style jsx>{`
        .spinner {
          display: inline-block;
          position: relative;
          width: 24px;
          height: 24px;
        }
        .spinner span {
          display: block;
          transform-origin: 12px 12px;
          animation: spinner 1.2s linear infinite;
        }
        .spinner span:after {
          content: " ";
          display: block;
          position: absolute;
          top: 0px;
          left: 11px;
          width: 2px;
          height: 5px;
          border-radius: 20%;
          background: currentColor;
        }
        .spinner span:nth-child(1) {
          transform: rotate(0deg);
          animation-delay: -1.1s;
        }
        .spinner span:nth-child(2) {
          transform: rotate(30deg);
          animation-delay: -1s;
        }
        .spinner span:nth-child(3) {
          transform: rotate(60deg);
          animation-delay: -0.9s;
        }
        .spinner span:nth-child(4) {
          transform: rotate(90deg);
          animation-delay: -0.8s;
        }
        .spinner span:nth-child(5) {
          transform: rotate(120deg);
          animation-delay: -0.7s;
        }
        .spinner span:nth-child(6) {
          transform: rotate(150deg);
          animation-delay: -0.6s;
        }
        .spinner span:nth-child(7) {
          transform: rotate(180deg);
          animation-delay: -0.5s;
        }
        .spinner span:nth-child(8) {
          transform: rotate(210deg);
          animation-delay: -0.4s;
        }
        .spinner span:nth-child(9) {
          transform: rotate(240deg);
          animation-delay: -0.3s;
        }
        .spinner span:nth-child(10) {
          transform: rotate(270deg);
          animation-delay: -0.2s;
        }
        .spinner span:nth-child(11) {
          transform: rotate(300deg);
          animation-delay: -0.1s;
        }
        .spinner span:nth-child(12) {
          transform: rotate(330deg);
          animation-delay: 0s;
        }
        @keyframes spinner {
          0% {
            opacity: 1;
          }
          100% {
            opacity: 0;
          }
        }
      `}</style>
    </span>
  );
}
