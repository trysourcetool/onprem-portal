import { createFileRoute } from '@tanstack/react-router';
import { Copy } from 'lucide-react';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useAuth } from '@/components/provider/auth-provider';
import { Button } from '@/components/ui/button';

export default function Index() {
  const { account } = useAuth();
  const onCopy = async (value: string) => {
    try {
      await navigator.clipboard.writeText(value);
      toast('License Key copied to clipboard');
    } catch (error) {
      console.error(error);
      toast('Failed to copy license key');
    }
  };
  return (
    <div>
      <div className={cn('flex flex-col gap-2')}>
        <h1 className="text-foreground text-3xl font-bold">License Key</h1>
      </div>
      <div className="flex w-screen flex-col gap-4 px-4 py-6 md:w-auto md:gap-6 md:px-6">
        <div className="w-full overflow-auto rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>License Key</TableHead>
                <TableHead className="w-[72px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell className="truncate">
                  {account?.license?.key}
                </TableCell>
                <TableCell align="center">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="cursor-pointer"
                    onClick={() => onCopy(account?.license?.key ?? '')}
                  >
                    <Copy className="size-4" />
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
          <div className="bg-muted border-t p-4"></div>
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_default/')({
  component: Index,
});
