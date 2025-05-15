import { Check } from 'lucide-react';
import clsx from 'clsx';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import type { FC } from 'react';
import { Card, CardContent } from '@/components/ui/card';

export const PricingCard: FC<{
  title: string;
  description: string;
  price: number;
  period?: string;
  buttonLabel: string;
  onClick: () => void;
  features: Array<string>;
  isPopular?: boolean;
  buttonDisabled?: boolean;
}> = ({
  title,
  description,
  price,
  period,
  buttonLabel,
  onClick,
  features,
  isPopular = false,
  buttonDisabled = false,
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
          <Button onClick={onClick} disabled={buttonDisabled}>
            {buttonLabel}
          </Button>
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
