import React from 'react'

const Button = ({
  children,
  onClick,
  secondary,
}) => (
  <button secondary={secondary ? 'secondary' : null} onClick={onClick}>
    {children}
    <style jsx>{`
      * {
        box-sizing: border-box;
        font-family: sans-serif;
      }

      button {
        border: none;
        font: inherit;
        display: inline-block;
        margin: 0;
        padding: 15px;
        width: 100%;
        color: white;
        height: 70px;
        border-radius: 6px;
        font-size: 28px;
        font-weight: 100;
        max-width: 350px;
        background: green;
      }

      button[secondary] {
        background: #666;
      }
    `}</style>
  </button>
)

export default Button
