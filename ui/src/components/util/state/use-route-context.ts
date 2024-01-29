import { BreadcrumbGroupProps } from '@cloudscape-design/components';
import { useMemo } from 'react';
import { useLocation } from 'react-router-dom';

interface RouteElement {
  path: string | RegExp;
  title?: string;
  breadcrumb?: string;
  children?: readonly RouteElement[];
}

const ROUTES = [{
  path: '',
  breadcrumb: 'Home',
  children: [],
}] satisfies readonly RouteElement[];

export function useRouteContext() {
  const location = useLocation();
  
  return useMemo(() => {
    let titlePrefix: string | undefined;
    const breadcrumbItems: BreadcrumbGroupProps.Item[] = [];

    if (location.pathname !== '/') {
      const parts = location.pathname.split('/').map(decodeURIComponent);

      let href = '';
      let routes: readonly RouteElement[] | undefined = ROUTES;

      for (const part of parts) {
        if (!href.endsWith('/')) {
          href += '/';
        }

        const currentHref = `${href}${encodeURIComponent(part)}`;
        let ignorePart = false;

        if (routes !== undefined) {
          let matchedRoute: RouteElement | undefined;

          for (const route of routes) {
            if ((route.path instanceof RegExp && route.path.test(part)) || route.path === part) {
              matchedRoute = route;
              break;
            }
          }

          if (matchedRoute !== undefined) {
            titlePrefix = matchedRoute.title;

            if (matchedRoute.breadcrumb !== undefined) {
              breadcrumbItems.push({
                text: matchedRoute.breadcrumb,
                href: currentHref,
              });
            }

            ignorePart = true;
          }

          routes = matchedRoute?.children;
        }

        if (!ignorePart) {
          breadcrumbItems.push({
            text: part,
            href: currentHref,
          });
        }

        href = currentHref;
      }
    }

    return {
      documentTitle: titlePrefix !== undefined ? `${titlePrefix} • GW2 Addon Manager` : 'GW2 Addon Manager',
      breadcrumbItems: breadcrumbItems,
    } as const;
  }, [location]);
}

export function useDocumentTitle() {
  const { documentTitle } = useRouteContext();
  return documentTitle;
}

export function useBreadcrumbItems() {
  const { breadcrumbItems } = useRouteContext();
  return breadcrumbItems;
}
