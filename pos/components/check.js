import React from 'react'

const Check = () => (
  <svg
    className="checkmark"
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 52 52"
  >
    <circle
      className="circle"
      cx="26"
      cy="26"
      r="25"
      fill="none"
    />
    <path
      className="check"
      fill="none"
      d="M14.1 27.2l7.1 7.2 16.7-16.8"
    />
    <style jsx>{`
      .circle {
        stroke-dasharray: 166;
        stroke-dashoffset: 166;
        stroke-width: 4;
        stroke-miterlimit: 10;
        stroke: green;
        fill: none;
        animation: stroke .6s cubic-bezier(0.650, 0.000, 0.450, 1.000) forwards;
      }

      .checkmark {
        width: 200px;
        height: 200px;
        border-radius: 50%;
        display: block;
        stroke-width: 4;
        stroke: #fff;
        stroke-miterlimit: 10;
        margin: 10% auto;
        box-shadow: inset 0px 0px 0px green;
        animation: fill .4s ease-in-out .4s forwards, scale .3s ease-in-out .9s both;
      }

      .check {
        transform-origin: 50% 50%;
        stroke-dasharray: 48;
        stroke-dashoffset: 48;
        animation: stroke .3s cubic-bezier(0.650, 0.000, 0.450, 1.000) .8s forwards;
      }

      @keyframes stroke {
        100% {
          stroke-dashoffset: 0;
        }
      }

      @keyframes scale {
        0%, 100% {
          transform: none;
        }
        50% {
          transform: scale3d(1.1, 1.1, 1);
        }
      }

      @keyframes fill {
        100% {
          box-shadow: inset 0px 0px 0px 100px green;
        }
      }
    `}</style>
  </svg>
)

export default Check
