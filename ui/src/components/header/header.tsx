import { TopNavigation, TopNavigationProps } from '@cloudscape-design/components';
import React, { useState } from 'react';
import { PreferencesModal } from '../preferences/preferences';
import { useAppControls } from '../util/context/app-controls';
import { useHttpClient } from '../util/context/http-client';
import { useDateFormat } from '../util/state/use-dateformat';
import classes from './header.module.scss';

export function Gw2AddonManagerHeader() {
  const { apiClient } = useHttpClient();
  const { notification } = useAppControls();
  const { formatDateTime } = useDateFormat();

  const [showPreferences, setShowPreferences] = useState(false);

  const utilities: TopNavigationProps.Utility[] = [
    {
      type: 'button',
      text: 'GitHub',
      href: 'https://github.com/its-felix/gw2-addon-manager',
      external: true,
      externalIconAriaLabel: '(opens in a new tab)',
    },
    {
      type: 'button',
      text: 'Preferences',
      iconName: 'settings',
      onClick: () => setShowPreferences(true),
    },
    {
      type: 'button',
      text: 'Shutdown',
      iconName: 'close',
      onClick: () => {},
    },
  ];

  return (
    <>
      <PreferencesModal visible={showPreferences} onDismiss={() => setShowPreferences(false)} />
      <header id="gw2am-custom-header" className={classes['gw2am-header']}>
        <TopNavigation
          identity={{
            href: '/',
            title: 'GW2 Addon Manager',
            logo: {
              src: '/logo_white.svg',
              alt: 'GW2 Addon Manager Logo',
            },
          }}
          utilities={utilities}
        />
      </header>
    </>
  );
}
