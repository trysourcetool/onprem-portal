import * as api from '@/api/instance';

export type Plan = {
  id: string;
  name: string;
  price: number;
  createdAt: string;
  updatedAt: string;
};

export const getPlans = async () => {
  const res = await api.get<{
    plans: Array<Plan>;
  }>({
    path: '/plans',
    auth: true,
  });

  return res;
};
