import * as api from '@/api/instance';

export type User = {
  createdAt: string;
  googleId: string;
  email: string;
  firstName: string;
  id: string;
  lastName: string;
  updatedAt: string;
};

export const getMe = async () => {
  const res = await api.get<{
    user: User;
  }>({ path: '/users/me', auth: true });

  return res;
};

export const updateMe = async (params: {
  data: {
    firstName?: string;
    lastName?: string;
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: '/users/me',
    data: params.data,
    auth: true,
  });

  return res;
};

export const sendUpdateMeEmailInstructions = async (params: {
  data: {
    email: string;
    emailConfirmation: string;
  };
}) => {
  const res = await api.post({
    path: '/users/me/email/instructions',
    data: params.data,
    auth: true,
  });

  return res;
};

export const updateMeEmail = async (params: {
  data: {
    token: string;
  };
}) => {
  const res = await api.put<{
    user: User;
  }>({
    path: '/users/me/email',
    data: params.data,
    auth: true,
  });

  return res;
};
