import React, {createContext, useContext, useEffect, useState,} from 'react';
import {BrowserStore} from '../../../lib/store/browser-store';
import {ConsentLevel} from '../../../lib/consent.model';

class ConsentOverrideBrowserStore extends BrowserStore {
  constructor(storage: Storage, ignoreClose: boolean, private readonly _hasConsent: boolean) {
    super(storage, ignoreClose);
  }

  override hasConsent(): boolean {
    return this._hasConsent;
  }
}

const DEFAULT = new ConsentOverrideBrowserStore(window.sessionStorage, true, true);
const BrowserStoreContext = createContext<BrowserStore>(DEFAULT);

export function BrowserStoreProvider({ storage, children }: React.PropsWithChildren<{ storage: Storage; }>) {
  const [store, setStore] = useState<BrowserStore>(DEFAULT);
  useEffect(() => {
    const v = new BrowserStore(storage, false);
    v.setConsentLevels([ConsentLevel.STRICTLY_NECESSARY, ConsentLevel.FUNCTIONALITY]);
    setStore(v);
    return () => v.close();
  }, [storage]);

  return (
    <BrowserStoreContext.Provider value={store}>
      {children}
    </BrowserStoreContext.Provider>
  );
}

export function useBrowserStore() {
  return useContext(BrowserStoreContext);
}
