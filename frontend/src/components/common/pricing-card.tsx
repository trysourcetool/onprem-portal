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
  onDialogSubmit?: () => void;
  features: Array<string>;
  isPopular?: boolean;
  buttonDisabled?: boolean;
  buttonType: 'docs' | 'dialog';
  isCurrentPlan?: boolean;
}> = ({
  title,
  description,
  price,
  period,
  onDialogSubmit,
  features,
  isPopular = false,
  buttonDisabled = false,
  buttonType,
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
          {buttonType === 'docs' && (
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
          {buttonType === 'dialog' && (
            <AlertDialog>
              <AlertDialogTrigger
                disabled={isCurrentPlan || buttonDisabled}
                asChild
              >
                <Button className="cursor-pointer">
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
                    onClick={onDialogSubmit}
                    className="cursor-pointer"
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
