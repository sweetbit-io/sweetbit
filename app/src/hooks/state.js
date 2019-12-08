import createPersistedState from 'use-persisted-state';

export const useDispenserState = createPersistedState('dispenser');
export const useNodesState = createPersistedState('nodes');
