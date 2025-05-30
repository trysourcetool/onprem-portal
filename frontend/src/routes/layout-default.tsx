import {
  Link,
  Outlet,
  createFileRoute,
  useLocation,
} from '@tanstack/react-router';
import { ChevronsUpDown, FileText, LogOut } from 'lucide-react';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar';
import { useAuth } from '@/components/provider/auth-provider';
import { Toaster } from '@/components/ui/sonner';
import { ModeToggle } from '@/components/common/mode-toggle';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

export default function DefaultLayout() {
  const { account, handleLogout } = useAuth();
  const { pathname } = useLocation();
  return (
    <SidebarProvider>
      {account && (
        <Sidebar collapsible="icon">
          <SidebarHeader>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground w-full cursor-default"
                >
                  <div className="flex flex-1 items-center gap-2 data-[state=open]:px-2 data-[state=open]:py-1">
                    <Link to={'/'} className="size-8">
                      <img
                        src="/images/logo-sidebar.png"
                        alt="Sourcetool"
                        className="size-full"
                      />
                    </Link>
                    <div className="flex flex-1 flex-col gap-0.5">
                      <p className="text-sidebar-foreground text-sm font-semibold">
                        Sourcetool
                      </p>
                    </div>
                    <ModeToggle />
                  </div>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarHeader>
          <SidebarContent>
            <SidebarGroup>
              <SidebarMenu>
                <SidebarMenuButton asChild isActive={pathname === '/'}>
                  <Link to={'/'}>
                    <FileText />
                    <span>License Key</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenu>
            </SidebarGroup>
          </SidebarContent>
          <SidebarFooter>
            <SidebarMenu>
              <SidebarMenuItem>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <SidebarMenuButton
                      size="lg"
                      className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                    >
                      <Avatar className="size-8 rounded-lg">
                        <AvatarFallback className="rounded-lg">
                          {account.firstName[0]}
                          {account.lastName[0]}
                        </AvatarFallback>
                      </Avatar>
                      <div className="grid flex-1 text-left text-sm leading-tight">
                        <span className="truncate font-semibold">
                          {account.firstName} {account.lastName}
                        </span>
                        <span className="truncate text-xs">
                          {account.email}
                        </span>
                      </div>
                      <ChevronsUpDown className="ml-auto size-4" />
                    </SidebarMenuButton>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                    side={'right'}
                    align="end"
                    sideOffset={4}
                  >
                    <DropdownMenuLabel className="p-0 font-normal">
                      <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                        <Avatar className="size-8 rounded-lg">
                          <AvatarFallback className="rounded-lg">
                            {account.firstName[0]}
                            {account.lastName[0]}
                          </AvatarFallback>
                        </Avatar>
                        <div className="grid flex-1 text-left text-sm leading-tight">
                          <span className="truncate font-semibold">
                            {account.firstName} {account.lastName}
                          </span>
                          <span className="truncate text-xs">
                            {account.email}
                          </span>
                        </div>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      onClick={handleLogout}
                      className="cursor-pointer"
                    >
                      <LogOut />
                      Log out
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarFooter>
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
          </div>
        </header>
        <main className="flex flex-1 flex-col px-4 py-6 md:px-6 md:py-6">
          <Outlet />
          <Toaster closeButton />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}

export const Route = createFileRoute('/_default')({
  component: DefaultLayout,
});
