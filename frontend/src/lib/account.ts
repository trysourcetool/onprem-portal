import dayjs from 'dayjs';
import type { User } from '@/api/modules/users';
import type { Subscription } from '@/api/modules/subscriptions';

const ACCOUNT_EXPIRATION_MONTHS = 1;

// Returns the remaining trial days and days until account deletion
export const checkAccountExpiredDays = (account: User) => {
  const now = dayjs();
  const createdAt = dayjs.unix(Number(account.createdAt));
  const expirationDate = createdAt.add(ACCOUNT_EXPIRATION_MONTHS, 'month');

  const expiredDays = dayjs(expirationDate).diff(now, 'day');

  return expiredDays < 0 ? 0 : expiredDays + 1;
};

export const checkTrialExpiredDays = (subscription: Subscription) => {
  if (subscription.status === 'trial') {
    return dayjs.unix(Number(subscription.trialEnd)).diff(dayjs(), 'day') + 1;
  }

  return 0;
};
