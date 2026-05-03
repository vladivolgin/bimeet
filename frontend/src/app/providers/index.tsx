import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ChakraProvider } from '@chakra-ui/react'
import theme from '@/app/styles/theme'
import { AuthProvider } from '@/features/auth/model/AuthContext'
import type { ReactNode } from 'react'
import '@/shared/i18n'

const queryClient = new QueryClient()

export default function Providers({ children }: { children: ReactNode }) {
  return (
    <ChakraProvider theme={theme}>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          {children}
        </AuthProvider>
      </QueryClientProvider>
    </ChakraProvider>
  )
}
