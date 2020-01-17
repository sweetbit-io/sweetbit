import React from 'react';
import css from 'styled-jsx/css';

export const { className: formClassName, styles: formStyles } = css.resolve`
  .group {
    position: relative;
    margin-top: 45px;
  }
  .group.centered input,
  .group.centered [contenteditable],
  .group.centered select {
    text-align: center;
  }
  .group input,
  .group textarea,
  .group [contenteditable],
  .group select {
    font-size: 18px;
    padding: 10px;
    display: block;
    width: 100%;
    border: none;
    border-bottom: 3px solid #ccc;
    background: transparent;
    border-radius: 0;
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
  }
  .group input:focus,
  .group textarea:focus,
  .group [contenteditable]:focus,
  .group select:focus {
    outline: none;
  }
  .group input::placeholder,
  .group textarea::placeholder,
  .group select::placeholder {
    color: transparent;
  }
  .group label {
    color: #999;
    font-size: 18px;
    font-weight: normal;
    position: absolute;
    pointer-events: none;
    left: 5px;
    top: 10px;
    transition: 0.2s ease all;
  }
  .group.centered label {
    left: 50%;
    transform: translateX(-50%);
  }
  .group input:focus ~ label,
  .group textarea:focus ~ label,
  .group [contenteditable]:focus ~ label,
  .group select:focus ~ label,
  .group input:not(:placeholder-shown) ~ label,
  .group textarea:not(:placeholder-shown) ~ label,
  .group input:valid ~ label,
  .group textarea:valid ~ label,
  .group [contenteditable]:not(:empty) ~ label,
  .group select:valid ~ label {
    top: -20px;
    font-size: 14px;
    color: #5264AE;
  }
`;
