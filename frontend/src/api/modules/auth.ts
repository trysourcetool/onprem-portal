import * as api from '@/api/instance';

export const logout = async () => {
  const res = await api.post({
    path: '/auth/logout',
    auth: true,
  });
  return res;
};

export const refreshToken = async () => {
  const res = await api.post<{
    expiresAt: string;
  }>({
    path: '/auth/refreshToken',
    auth: true,
  });

  return res;
};

export const requestGoogleAuthLink = async () => {
  const res = await api.post<{
    authUrl: string;
  }>({
    path: '/auth/google/request',
  });

  return res;
};

export const authenticateWithGoogle = async (params: {
  data: {
    code: string;
    state: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
    isNewUser: boolean;
  }>({
    path: '/auth/google/authenticate',
    data: params.data,
  });

  return res;
};

export const registerWithGoogle = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.post<{
    authUrl: string;
    token: string;
  }>({
    path: '/auth/google/register',
    data: params.data,
  });

  return res;
};

export const requestMagicLink = async (params: {
  data: {
    email: string;
  };
}) => {
  const res = await api.post<{
    email: string;
    isNew: boolean;
  }>({
    path: '/auth/magic/request',
    data: params.data,
  });

  return res;
};

export const authenticateWithMagicLink = async (params: {
  data: {
    token: string;
    firstName?: string;
    lastName?: string;
  };
}) => {
  const res = await api.post<{
    registrationToken: string;
    expiresAt: string;
    isNewUser: boolean;
  }>({
    path: '/auth/magic/authenticate',
    data: params.data,
  });

  return res;
};

export const registerWithMagicLink = async (params: {
  data: {
    token: string;
    firstName: string;
    lastName: string;
  };
}) => {
  const res = await api.post<{
    hasOrganization: boolean;
    expiresAt: string;
  }>({
    path: '/auth/magic/register',
    data: params.data,
  });

  return res;
};
