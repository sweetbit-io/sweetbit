import { useMemo } from 'react';

export function useApi({
  publicUrl,
  publicWsUrl,
}) {
  return useMemo(() => ({
    async addNode(data) {
      const res = await fetch(`${publicUrl}/api/v1/nodes`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      return res.json();
    },
    async renameNode(id, name) {
      const res = await fetch(`${publicUrl}/api/v1/nodes/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          op: 'rename',
          name,
        }),
      });
      return await res.json();
    },
    async enableNode(id, enabled) {
      const res = await fetch(`${publicUrl}/api/v1/nodes/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          op: 'enable',
          enabled,
        }),
      });
      return await res.json();
    },
    async unlockNode(id, password) {
      const res = await fetch(`${publicUrl}/api/v1/nodes/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          op: 'unlock',
          password,
        }),
      });
      return await res.json();
    },
    async initNode(id, password, mnemonic) {
      const res = await fetch(`${publicUrl}/api/v1/nodes/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          op: 'init',
          password,
          mnemonic,
        }),
      });

      if (res.status !== 200) {
        const { error } = await res.json();
        throw new Error(error);
      }

      return await res.json();
    },
    subscribeNodeStatus(id) {
      return new WebSocket(`${publicWsUrl}/api/v1/nodes/${id}/status`);
    },
    async generateNodeSeed(id) {
      const res = await fetch(`${publicUrl}/api/v1/nodes/${id}/seed`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
        }),
      });

      if (res.status !== 200) {
        const { error } = await res.json();
        throw new Error(error);
      }

      return res.json();
    },
    async generateNodeConnection(id) {
      const res = await fetch(`${publicUrl}/api/v1/nodes/${id}/connection`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
        }),
      });

      if (res.status !== 200) {
        const { error } = await res.json();
        throw new Error(error);
      }

      return res.json();
    },
  }), [publicUrl, publicWsUrl]);
}
