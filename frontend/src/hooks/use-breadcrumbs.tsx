import { createContext, useContext, useState } from 'react';
import type { Dispatch, ReactNode, SetStateAction } from 'react';

type Breadcrumb = {
  to?: string;
  label: string;
};

type BreadcrumbsState = {
  breadcrumbsState?: Array<Breadcrumb>;
  setBreadcrumbsState?: Dispatch<SetStateAction<Array<Breadcrumb>>>;
};

export const breadcrumbsContext = createContext<BreadcrumbsState>({});

export function BreadcrumbsProvider(props: { children: ReactNode }) {
  const [breadcrumbsState, setBreadcrumbsState] = useState<Array<Breadcrumb>>(
    [],
  );

  return (
    <breadcrumbsContext.Provider
      value={{ breadcrumbsState, setBreadcrumbsState }}
    >
      {props.children}
    </breadcrumbsContext.Provider>
  );
}

export const useBreadcrumbs = () => {
  const { breadcrumbsState, setBreadcrumbsState } =
    useContext(breadcrumbsContext);

  return { breadcrumbsState, setBreadcrumbsState };
};
