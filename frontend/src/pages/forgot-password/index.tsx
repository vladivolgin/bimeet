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
} from '@chakra-ui/react'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { FiMoon, FiSun } from 'react-icons/fi'
import { authApi } from '@/features/auth/api'

export default function ForgotPasswordPage() {
  const { t } = useTranslation()
  const { colorMode, toggleColorMode } = useColorMode()
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [submitted, setSubmitted] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!email) return
    setLoading(true)
    try {
      await authApi.forgotPassword(email)
    } catch {
      // Always show "sent" — don't leak whether email exists
    }
    setLoading(false)
    setSubmitted(true)
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
              {submitted ? (
                <VStack spacing={5} align="stretch">
                  <Box>
                    <Heading size="md" mb={1}>{t('auth.resetEmailSent')}</Heading>
                    <Text fontSize="sm" color="dimText">{t('auth.resetEmailSentDescription')}</Text>
                  </Box>
                  <Button as={Link} to="/login" colorScheme="blue" w="full" size="md">
                    {t('auth.backToLogin')}
                  </Button>
                </VStack>
              ) : (
                <VStack spacing={5} as="form" onSubmit={handleSubmit}>
                  <Box w="full">
                    <Heading size="md" mb={1}>{t('auth.forgotPasswordTitle')}</Heading>
                    <Text fontSize="sm" color="dimText">{t('auth.forgotPasswordSubtitle')}</Text>
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

                  <Button type="submit" colorScheme="blue" w="full" size="md" isLoading={loading} mt={1}>
                    {t('auth.sendResetLink')}
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
              )}
            </Box>
          </VStack>
        </Box>
      </Box>
    </Box>
  )
}
