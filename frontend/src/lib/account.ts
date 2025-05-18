import dayjs from 'dayjs';
import type { User } from '@/api/modules/users';
import type { Subscription } from '@/api/modules/subscriptions';

export const getAccountDeletionRemainingDays = (account: User | null) => {
  if (!account) return 0;

  const now = dayjs();
  const scheduledDeletionAt = dayjs.unix(Number(account.scheduledDeletionAt));

  const expiredDays = scheduledDeletionAt.diff(now, 'day');

  return expiredDays < 0 ? 0 : expiredDays + 1;
};

export const getTrialStatus = (subscription: Subscription | null) => {
  if (
    subscription &&
    subscription?.status === 'trial' &&
    subscription.trialEnd
  ) {
    const trialEndDay = dayjs.unix(Number(subscription.trialEnd));
    return {
      expiredDays: trialEndDay.diff(dayjs(), 'day') + 1,
      isTrialEnd: trialEndDay.isBefore(dayjs()),
    };
  }

  return {
    expiredDays: 0,
    isTrialEnd: false,
  };
};
