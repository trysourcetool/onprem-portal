import * as users from './modules/users';
import * as auth from './modules/auth';
import * as subscriptions from './modules/subscriptions';
import * as stripe from './modules/stripe';
import * as plan from './modules/plan';

export const api = {
  auth,
  users,
  subscriptions,
  stripe,
  plan,
};
