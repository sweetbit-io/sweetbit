import React, { useRef, useCallback } from 'react';

export default function AddNode({
  onCancel,
  onAdd,
}) {
  const inputEl = useRef(null);
  const submit = useCallback((e) => {
    e.preventDefault();
    onAdd(inputEl.current.value);
  }, [inputEl, onAdd]);

  return (
    <div>
      <span>The count is</span>
      <form onSubmit={submit}>
        <input type="text" ref={inputEl} />
        <button type="submit">add</button>
        <button type="button" onClick={onCancel}>cancel</button>
      </form>
      <style jsx>{`
        div {
          padding: 20px;
        }
      `}</style>
    </div>
  );
}
