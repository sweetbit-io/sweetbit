import React from 'react';

export default function Mnemonic({
  words,
  onConfirm,
}) {
  return (
    <div>
      {words.map((word) => (
        <div key={word}>{word}</div>
      ))}
      <style jsx>{`
      `}</style>
    </div>
  );
}
