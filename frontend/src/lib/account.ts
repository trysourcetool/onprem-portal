import dayjs from 'dayjs';
import type { User } from '@/api/modules/users';

const ACCOUNT_EXPIRATION_MONTHS = 1;
const TRIAL_EXPIRATION_DAYS = 14;

// Returns the remaining trial days and days until account deletion
export const checkAccountExpiredDays = (account: User) => {
  const now = dayjs();
  const createdAt = dayjs.unix(Number(account.createdAt));
  const expirationDate = createdAt.add(ACCOUNT_EXPIRATION_MONTHS, 'month');
  const trialExpirationDate = createdAt.add(TRIAL_EXPIRATION_DAYS, 'day');

  const expiredDays = dayjs(expirationDate).diff(now, 'day');
  const trialExpiredDays = dayjs(trialExpirationDate).diff(now, 'day');

  return {
    expiredDays: expiredDays < 0 ? 0 : expiredDays + 1,
    trialExpiredDays: trialExpiredDays < 0 ? 0 : trialExpiredDays + 1,
  };
};
