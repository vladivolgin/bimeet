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
  HStack,
  IconButton,
  Tooltip,
  useColorMode,
  useToast,
} from '@chakra-ui/react'
import { useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { FiMoon, FiSun } from 'react-icons/fi'
import { authApi } from '@/features/auth/api'

export default function ResetPasswordPage() {
  const { t } = useTranslation()
  const toast = useToast()
  const navigate = useNavigate()
  const { colorMode, toggleColorMode } = useColorMode()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') ?? ''
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!token) {
      toast({ title: t('auth.missingResetToken'), status: 'error', duration: 3000 })
      return
    }
    if (!password || !confirm) {
      toast({ title: t('auth.fillAllFields'), status: 'warning', duration: 3000 })
      return
    }
    if (password !== confirm) {
      toast({ title: t('auth.passwordsDoNotMatch'), status: 'error', duration: 3000 })
      return
    }
    setLoading(true)
    try {
      await authApi.resetPassword(token, password)
      toast({ title: t('auth.resetPasswordSuccess'), status: 'success', duration: 3000 })
      navigate('/login')
    } catch {
      toast({ title: t('auth.resetPasswordError'), status: 'error', duration: 4000 })
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
                  <Heading size="md" mb={1}>{t('auth.resetPasswordTitle')}</Heading>
                  <Text fontSize="sm" color="dimText">{t('auth.resetPasswordSubtitle')}</Text>
                </Box>

                <FormControl>
                  <FormLabel fontSize="sm" fontWeight="600">{t('auth.newPassword')}</FormLabel>
                  <Input
                    type="password"
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                    placeholder={t('auth.passwordPlaceholder')}
                    size="md"
                  />
                </FormControl>

                <FormControl>
                  <FormLabel fontSize="sm" fontWeight="600">{t('auth.confirmPassword')}</FormLabel>
                  <Input
                    type="password"
                    value={confirm}
                    onChange={e => setConfirm(e.target.value)}
                    placeholder={t('auth.passwordPlaceholder')}
                    size="md"
                  />
                </FormControl>

                <Button type="submit" colorScheme="blue" w="full" size="md" isLoading={loading} mt={1}>
                  {t('auth.resetPassword')}
                </Button>

                <Text fontSize="sm" color="dimText" textAlign="center">
                  <ChakraLink
                    as={Link}
                    to="/login"
                    color="brand.600"
                    fontWeight="600"
                    _hover={{ textDecoration: 'none', color: 'brand.700' }}
                  >
                    {t('auth.backToLogin')}
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
