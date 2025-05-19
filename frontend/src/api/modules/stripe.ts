import * as api from '@/api/instance';

export const createCheckoutSession = async (params: {
  data: {
    planId: string;
  };
}) => {
  const res = await api.post<{
    url: string;
  }>({
    path: '/stripe/createCheckoutSession',
    data: params.data,
    auth: true,
  });

  return res;
};

export const getCustomerPortalUrl = async () => {
  const res = await api.get<{
    url: string;
  }>({ path: '/stripe/customerPortalUrl', auth: true });

  return res;
};
