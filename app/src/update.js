import React from 'react';
import css from 'styled-jsx/css';
import Markdown from 'react-markdown';
import Button from './button';
import { ReactComponent as UpdateImage } from './update.svg';

const { className, styles } = css.resolve`
  .image {
    width: auto;
    height: 120px;
  }
`;

export default function Update({
  onCancel,
  onUpdate,
  name,
  body,
}) {
  return (
    <div className="update">
      <p className="center">
        <UpdateImage className={`${className} image`} />
      </p>
      <h1 className="center">{name}</h1>
      <Markdown source={body} />
      <div className="center actions">
        <Button type="button" onClick={onUpdate}>update</Button>
        <span> </span>
        <Button type="button" onClick={onCancel} outline>cancel</Button>
      </div>
      {styles}
      <style jsx>{`
        .update {
          padding: 20px;
        }

        .center {
          text-align: center;
        }

        .actions {
          padding-top: 40px;
        }

        h1 {
          margin: 0;
        }

        .update :global(img) {
          max-width: 100%;
          width: 100%;
          height: auto;
        }
      `}</style>
    </div>
  );
}
