import {
  AppLayout,
  AppLayoutProps,
  Flashbar,
  FlashbarProps,
  NonCancelableCustomEvent,
} from '@cloudscape-design/components';
import { I18nProvider as CSI18nProvider } from '@cloudscape-design/components/i18n';
import enMessages from '@cloudscape-design/components/i18n/messages/all.en';
import {
  applyDensity, applyMode, Density, Mode, 
} from '@cloudscape-design/global-styles';
import React, {
  createContext, useContext, useEffect, useMemo, useState, 
} from 'react';
import { ColorScheme, UIDensity } from '../lib/preferences.model';
import { Breadcrumb } from './breadcrumb/breadcrumb';
import { Gw2AddonManagerFooter } from './footer/footer';
import { Gw2AddonManagerHeader } from './header/header';
import { AppControlsProvider } from './util/context/app-controls';
import { BrowserStoreProvider } from './util/context/browser-store';
import { HttpClientProvider } from './util/context/http-client';
import { usePreferences } from './util/state/use-preferences';
import { useDocumentTitle } from './util/state/use-route-context';

interface AppControlsState {
  tools: {
    element: React.ReactNode | undefined;
    open: boolean;
    onChange: (e: NonCancelableCustomEvent<AppLayoutProps.ChangeDetail>) => void;
  };
  notification: {
    messages: Array<FlashbarProps.MessageDefinition>;
  };
}

const AppControlsStateContext = createContext<AppControlsState>({
  tools: {
    element: undefined,
    open: false,
    onChange: () => {},
  },
  notification: {
    messages: [],
  },
});

export interface RootLayoutProps extends Omit<AppLayoutProps, 'content'> {
  headerHide: boolean;
  breadcrumbsHide: boolean;
}

export function RootLayout({
  headerHide, breadcrumbsHide, children, ...appLayoutProps 
}: React.PropsWithChildren<RootLayoutProps>) {
  const documentTitle = useDocumentTitle();
  const appControlsState = useContext(AppControlsStateContext);

  useEffect(() => {
    const restore = document.title;
    document.title = documentTitle;
    return () => { document.title = restore; };
  }, [documentTitle]);

  return (
    <>
      {!headerHide && <Gw2AddonManagerHeader />}
      <HeaderSelectorFixAppLayout
        toolsHide={appControlsState.tools.element === undefined}
        tools={appControlsState.tools.element}
        toolsOpen={appControlsState.tools.element !== undefined && appControlsState.tools.open}
        onToolsChange={appControlsState.tools.onChange}
        headerSelector={headerHide ? undefined : '#gw2am-custom-header'}
        footerSelector={'#gw2am-custom-footer'}
        stickyNotifications={true}
        notifications={<Flashbar stackItems={true} items={appControlsState.notification.messages} />}
        breadcrumbs={breadcrumbsHide ? undefined : <Breadcrumb />}
        navigationHide={true}
        content={children}
        {...appLayoutProps}
      />
      <Gw2AddonManagerFooter />
    </>
  );
}

function HeaderSelectorFixAppLayout(props: AppLayoutProps) {
  const { headerSelector, ...appLayoutProps } = props;
  const [key, setKey] = useState(`a${Date.now()}-${Math.random()}`);

  useEffect(() => {
    setKey(`a${Date.now()}-${Math.random()}`);
  }, [headerSelector]);

  return (
    <AppLayout key={key} headerSelector={headerSelector} {...appLayoutProps} />
  );
}

export function BaseProviders({ children }: React.PropsWithChildren) {
  return (
    <BrowserStoreProvider storage={window.localStorage}>
      <HttpClientProvider>
        <InternalBaseProviders>
          {children}
        </InternalBaseProviders>
      </HttpClientProvider>
    </BrowserStoreProvider>
  );
}

function InternalBaseProviders({ children }: React.PropsWithChildren) {
  const [preferences] = usePreferences();
  const [tools, setTools] = useState<React.ReactNode>();
  const [toolsOpen, setToolsOpen] = useState(false);
  const [notificationMessages, setNotificationMessages] = useState<Array<FlashbarProps.MessageDefinition>>([]);

  useEffect(() => {
    document.getElementById('temp_style')?.remove();
  }, []);

  useEffect(() => {
    applyMode(preferences.effectiveColorScheme === ColorScheme.LIGHT ? Mode.Light : Mode.Dark);
    applyDensity(preferences.uiDensity === UIDensity.COMFORTABLE ? Density.Comfortable : Density.Compact);
  }, [preferences]);

  const appControlsState = useMemo<AppControlsState>(() => ({
    tools: {
      element: tools,
      open: toolsOpen,
      onChange(e): void {
        setToolsOpen(e.detail.open);
      },
    },
    notification: {
      messages: notificationMessages,
    },
  }), [tools, toolsOpen, notificationMessages]);

  return (
    <CSI18nProvider locale={'en'} messages={[enMessages]}>
      <AppControlsProvider setTools={setTools} setToolsOpen={setToolsOpen} setNotificationMessages={setNotificationMessages}>
        <AppControlsStateContext.Provider value={appControlsState}>
          {children}
        </AppControlsStateContext.Provider>
      </AppControlsProvider>
    </CSI18nProvider>
  );
}
