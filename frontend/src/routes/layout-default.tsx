import { Outlet, createFileRoute } from '@tanstack/react-router';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar';
import { useAuth } from '@/components/provider/auth-provider';
import { Toaster } from '@/components/ui/sonner';
import { ModeToggle } from '@/components/common/mode-toggle';

export default function DefaultLayout() {
  const { account } = useAuth();

  return (
    <SidebarProvider>
      {account && (
        <Sidebar collapsible="icon">
          <SidebarHeader />
          <SidebarContent></SidebarContent>
          <SidebarFooter />
        </Sidebar>
      )}
      <SidebarInset className="flex flex-col">
        <header className="bg-background group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 sticky top-0 z-10 flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
          <div className="flex flex-1 items-center gap-2 px-4">
            <div className="flex-1">
              {account ? (
                <SidebarTrigger className="-ml-1" />
              ) : (
                <div className="size-8">
                  <img
                    src="/images/logo.png"
                    alt="Sourcetool"
                    className="size-full"
                  />
                </div>
              )}
            </div>
            <ModeToggle />
          </div>
        </header>
        <main className="flex flex-1 flex-col px-4 py-6 md:px-6 md:py-6">
          <Outlet />
          <Toaster />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}

export const Route = createFileRoute('/_default')({
  component: DefaultLayout,
});
