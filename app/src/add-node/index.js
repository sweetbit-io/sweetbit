import React, { useState, useCallback } from 'react';
import SelectType from './select-type';
import RemoteLnd from './remote-lnd';

export default function AddNode({
  api,
  onCancel,
  dispatchNodesAction,
}) {
  const [state, setState] = useState({
    type: null,
  });

  const selectType = useCallback((type) => {
    const onSelect = async () => {
      if (type === 'local') {
        let node = await api.addNode({
          type: 'local',
          name: 'Local Node',
        });

        node = await api.enableNode(node.id, true);

        dispatchNodesAction({
          type: 'add',
          node,
        });

        onCancel();
      } else {
        setState({ type });
      }
    };
    onSelect();
  }, [api, dispatchNodesAction, onCancel])

  return (
    <div className="add">
      {state.type === null ? (
        <SelectType onSelect={selectType} onCancel={onCancel} />
      ) : state.type === 'local' ? (
        <RemoteLnd onCancel={onCancel} />
      ) : (
        <RemoteLnd onCancel={onCancel} />
      )}
    </div>
  );
}
