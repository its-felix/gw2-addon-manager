import { applyMode, Mode } from '@cloudscape-design/global-styles';
import React from 'react';
import ReactDOM from 'react-dom/client';
import {
  createBrowserRouter, Outlet, RouterProvider,
} from 'react-router-dom';
import '@cloudscape-design/global-styles/index.css';
import { BaseProviders, RootLayout } from '../components/root';
import { ErrorPage } from './error-page';
import { Home } from './home';

// region router
const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <RootLayout headerHide={false} breadcrumbsHide={false}>
        <Outlet />
      </RootLayout>
    ),
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <Home /> },
    ],
  },
]);
// endregion

const root = ReactDOM.createRoot(document.getElementById('root')!);
const element = (
  <React.StrictMode>
    <BaseProviders>
      <RouterProvider router={router} />
    </BaseProviders>
  </React.StrictMode>
);

applyMode(Mode.Dark);
root.render(element);
