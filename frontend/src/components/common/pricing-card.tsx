import { Check } from 'lucide-react';
import clsx from 'clsx';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '../ui/alert-dialog';
import type { FC } from 'react';
import { Card, CardContent } from '@/components/ui/card';

export const PricingCard: FC<{
  title: string;
  description: string;
  price: number;
  period?: string;
  onUpgrade?: () => void;
  features: Array<string>;
  isPopular?: boolean;
  buttonDisabled?: boolean;
  planType: 'free' | 'team' | 'business';
  hasPlan?: boolean;
  isCurrentPlan?: boolean;
}> = ({
  title,
  description,
  price,
  period,
  onUpgrade,
  features,
  isPopular = false,
  buttonDisabled = false,
  planType,
  hasPlan = false,
  isCurrentPlan = false,
}) => {
  return (
    <Card
      className={clsx('p-6 lg:p-8', isPopular && 'border-primary border-2')}
    >
      <CardContent className="flex flex-col gap-8 p-0">
        <div className="flex flex-col gap-6">
          <div className="relative flex flex-col gap-3">
            {isPopular && (
              <Badge className="absolute right-0 top-1">Most popular</Badge>
            )}
            <h3 className="text-lg font-semibold">{title}</h3>
            <p className="text-muted-foreground text-sm">{description}</p>
          </div>
          <div className="flex items-end gap-0.5">
            <span className="text-4xl font-semibold">${price}</span>
            {period && (
              <span className="text-muted-foreground text-base">/{period}</span>
            )}
          </div>
          {planType === 'free' && (
            <Button
              disabled={buttonDisabled}
              asChild
              className="cursor-pointer"
            >
              <a target="_blank" rel="noreferrer">
                Read documentation
              </a>
            </Button>
          )}
          {(planType === 'team' || planType === 'business') && !hasPlan && (
            <Button
              disabled={buttonDisabled}
              className="cursor-pointer"
              onClick={onUpgrade}
            >
              Upgrade
            </Button>
          )}
          {(planType === 'team' || planType === 'business') && hasPlan && (
            <AlertDialog>
              <AlertDialogTrigger
                disabled={isCurrentPlan || buttonDisabled}
                asChild
              >
                <Button
                  className="cursor-pointer"
                  variant={isCurrentPlan ? 'secondary' : 'default'}
                >
                  {isCurrentPlan ? 'Current Plan' : 'Upgrade'}
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Upgrade your plan</AlertDialogTitle>
                  <AlertDialogDescription>
                    You are about to upgrade to Business plan.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel className="cursor-pointer">
                    Cancel
                  </AlertDialogCancel>
                  <AlertDialogAction
                    onClick={() => !isCurrentPlan && onUpgrade?.()}
                    className="cursor-pointer"
                    disabled={isCurrentPlan}
                  >
                    Upgrade
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          )}
        </div>
        <div className="flex flex-col gap-4">
          {features.map((feature) => (
            <div key={feature} className="flex items-center gap-3">
              <Check className="size-5" />
              <span className="text-muted-foreground flex-1 text-sm">
                {feature}
              </span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};
