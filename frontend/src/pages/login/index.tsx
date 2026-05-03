import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
  Heading,
  Text,
  Link as ChakraLink,
  useToast,
  HStack,
  IconButton,
  Tooltip,
  useColorMode,
} from '@chakra-ui/react'
import { useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { FiMoon, FiSun } from 'react-icons/fi'
import { useAuth } from '@/features/auth/model/AuthContext'
import { apiFetch } from '@/shared/api/client'

export default function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const toast = useToast()
  const { t } = useTranslation()
  const { colorMode, toggleColorMode } = useColorMode()
  const [searchParams] = useSearchParams()
  const inviteToken = searchParams.get('token')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    const ok = await login(email, password)
    if (!ok) {
      setLoading(false)
      toast({ title: t('auth.loginError'), status: 'error', duration: 3000 })
      return
    }
    if (inviteToken) {
      try {
        const participant = await apiFetch<{ event_id: string }>(`/api/events/invite/${inviteToken}`, {
          method: 'POST',
          body: JSON.stringify({ action: 'join' }),
        })
        navigate(`/events/${participant.event_id}`)
      } catch {
        navigate('/events')
      }
    } else {
      navigate('/events')
    }
    setLoading(false)
  }

  return (
    <Box minH="100vh" display="flex" position="relative" overflow="hidden">
      <Box position="absolute" inset={0} bg="pageBg" />
      <Box
        position="absolute"
        top="-160px"
        left="-160px"
        w="480px"
        h="480px"
        borderRadius="full"
        bgGradient="radial(#C7D2FE, transparent)"
        opacity={0.5}
        filter="blur(80px)"
      />
      <Box
        position="absolute"
        bottom="-140px"
        right="-140px"
        w="420px"
        h="420px"
        borderRadius="full"
        bgGradient="radial(#A5B4FC, transparent)"
        opacity={0.3}
        filter="blur(100px)"
      />

      <Tooltip label={colorMode === 'light' ? t('layout.darkTheme') : t('layout.lightTheme')} borderRadius="lg">
        <IconButton
          aria-label={t('layout.toggleTheme')}
          icon={colorMode === 'light' ? <FiMoon size={15} /> : <FiSun size={15} />}
          size="sm"
          variant="ghost"
          position="absolute"
          top={4}
          right={4}
          zIndex={2}
          onClick={toggleColorMode}
        />
      </Tooltip>

      <Box
        position="relative"
        zIndex={1}
        w="full"
        display="flex"
        alignItems="center"
        justifyContent="center"
        p={4}
      >
        <Box w="full" maxW="380px">
          <VStack spacing={7} align="stretch">
            <VStack spacing={1.5} align="center">
              <HStack spacing={2}>
                <Box
                  w="28px" h="28px" borderRadius="md" bg="brand.600"
                  display="flex" alignItems="center" justifyContent="center"
                >
                  <Text fontSize="12px" lineHeight={1} color="white">✦</Text>
                </Box>
                <Text fontSize="lg" fontWeight="600" color="mainText" letterSpacing="-0.3px">Bimeet</Text>
              </HStack>
              <Text color="dimText" fontSize="sm">{t('auth.tagline')}</Text>
            </VStack>

            <Box
              bg="cardBg"
              borderRadius="xl"
              p={7}
              boxShadow="0 2px 16px rgba(0,0,0,0.07), 0 1px 3px rgba(0,0,0,0.04)"
              border="1px solid"
              borderColor="cardBorder"
            >
              <VStack spacing={5} as="form" onSubmit={handleSubmit}>
                <Box w="full">
                  <Heading size="md" mb={1}>{t('auth.welcome')}</Heading>
                  <Text fontSize="sm" color="dimText">{t('auth.signInSubtitle')}</Text>
                </Box>

                <FormControl>
                  <FormLabel fontSize="sm" fontWeight="600">{t('auth.email')}</FormLabel>
                  <Input
                    type="email"
                    value={email}
                    onChange={e => setEmail(e.target.value)}
                    placeholder={t('auth.emailPlaceholder')}
                    size="md"
                  />
                </FormControl>

                <FormControl>
                  <FormLabel fontSize="sm" fontWeight="600">{t('auth.password')}</FormLabel>
                  <Input
                    type="password"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                    placeholder={t('auth.passwordPlaceholder')}
                    size="md"
                  />
                </FormControl>

                <Button type="submit" colorScheme="blue" w="full" size="md" isLoading={loading} mt={1}>
                  {t('auth.signIn')}
                </Button>

                <Text fontSize="sm" textAlign="center">
                  <ChakraLink
                    as={Link}
                    to="/forgot-password"
                    color="brand.600"
                    fontWeight="600"
                    _hover={{ textDecoration: 'none', color: 'brand.700' }}
                  >
                    {t('auth.forgotPassword')}
                  </ChakraLink>
                </Text>

                <Text fontSize="sm" color="dimText" textAlign="center">
                  {t('auth.noAccount')}{' '}
                  <ChakraLink
                    as={Link}
                    to={inviteToken ? `/register?token=${inviteToken}` : '/register'}
                    color="brand.600"
                    fontWeight="600"
                    _hover={{ textDecoration: 'none', color: 'brand.700' }}
                  >
                    {t('auth.signUp')}
                  </ChakraLink>
                </Text>
              </VStack>
            </Box>

          </VStack>
        </Box>
      </Box>
    </Box>
  )
}
