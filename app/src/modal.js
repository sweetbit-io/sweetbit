import React from 'react';
import ReactModal from 'react-modal';
import css from 'styled-jsx/css';

ReactModal.setAppElement('#root')

const { className, styles } = css.resolve`
  .ReactModal__Overlay {
    position: fixed;
    overflow: scroll;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
  }
  .ReactModal__Overlay {
    opacity: 0;
    transition: opacity 200ms ease-in-out;
  }
  .ReactModal__Overlay--after-open {
    opacity: 1;
  }
  .ReactModal__Overlay--before-close {
    opacity: 0;
  }
  .ReactModal__Content {
    display: block;
    box-shadow: 13px 13px 13px aliceblue;
    padding: 40px 0;
    max-width: 480px;
    margin: 0 auto;
  }
`

function Modal({ children, open, onExited }) {
  return (
    <ReactModal
      closeTimeoutMS={200}
      isOpen={open}
      className={className}
      overlayClassName={className}
      contentLabel="modal"
      onRequestClose={onExited}
    >
      {children}
      {styles}
    </ReactModal>
  );
}

export default Modal;
