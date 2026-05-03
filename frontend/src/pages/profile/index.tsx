import {
  Avatar,
  Box,
  Button,
  Divider,
  FormControl,
  FormLabel,
  Grid,
  Heading,
  HStack,
  Icon,
  Input,
  Spinner,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  useToast,
  VStack,
} from '@chakra-ui/react'
import { useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { FiCalendar, FiCamera, FiMail, FiMapPin, FiTrash2, FiUser } from 'react-icons/fi'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { authApi } from '@/features/auth/api'
import { useAuth } from '@/features/auth/model/AuthContext'
import { formatDate } from '@/shared/lib/formatDate'

export default function ProfilePage() {
  const { user, updateUser } = useAuth()
  const { t } = useTranslation()
  const toast = useToast()
  const queryClient = useQueryClient()
  const avatarInputRef = useRef<HTMLInputElement>(null)

  const { data: stats } = useQuery({
    queryKey: ['profile', 'stats'],
    queryFn: authApi.getStats,
  })

  const [form, setForm] = useState({
    name: user?.name ?? '',
    last_name: user?.last_name ?? '',
    city: user?.city ?? '',
    birth_date: user?.birth_date ? user.birth_date.slice(0, 10) : '',
  })

  const avatarMutation = useMutation({
    mutationFn: authApi.uploadAvatar,
    onSuccess: (updated) => {
      updateUser(updated)
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      toast({ title: t('profile.avatarSaved'), status: 'success', duration: 2500, isClosable: true })
    },
    onError: () => {
      toast({ title: t('profile.uploadError'), status: 'error', duration: 3000, isClosable: true })
    },
  })

  const deleteAvatarMutation = useMutation({
    mutationFn: authApi.deleteAvatar,
    onSuccess: (updated) => {
      updateUser(updated)
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      toast({ title: t('profile.avatarDeleted'), status: 'success', duration: 2500, isClosable: true })
    },
    onError: () => {
      toast({ title: t('profile.uploadError'), status: 'error', duration: 3000, isClosable: true })
    },
  })

  const updateMutation = useMutation({
    mutationFn: authApi.updateProfile,
    onSuccess: (updated) => {
      updateUser(updated)
      queryClient.invalidateQueries({ queryKey: ['profile'] })
      setForm({
        name: updated.name,
        last_name: updated.last_name ?? '',
        city: updated.city ?? '',
        birth_date: updated.birth_date ? updated.birth_date.slice(0, 10) : '',
      })
      toast({
        title: t('profile.saved'),
        status: 'success',
        duration: 2500,
        isClosable: true,
      })
    },
  })

  function handleSave() {
    updateMutation.mutate(form)
  }

  const fullName = [user?.name, user?.last_name].filter(Boolean).join(' ')

  return (
    <Grid
      templateColumns={{ base: '1fr', lg: '280px 1fr' }}
      gap={5}
      alignItems="flex-start"
    >
      {/* Left panel */}
      <Box position={{ lg: 'sticky' }} top={{ lg: '64px' }}>
        <Box
          bg="cardBg"
          borderRadius="xl"
          border="1px solid"
          borderColor="cardBorder"
          boxShadow="0 1px 3px rgba(15,23,42,0.04), 0 4px 16px rgba(15,23,42,0.05)"
          p={6}
        >
          <VStack spacing={4} align="center">
            <Box position="relative">
              <Box
                cursor="pointer"
                role="group"
                onClick={() => avatarInputRef.current?.click()}
              >
                <Avatar
                  size="2xl"
                  src={user?.avatar_url ?? undefined}
                  name={fullName || user?.name}
                  bg="brand.600"
                  color="white"
                  fontWeight="700"
                />
                <Box
                  position="absolute"
                  inset={0}
                  borderRadius="full"
                  bg="blackAlpha.600"
                  opacity={0}
                  _groupHover={{ opacity: 1 }}
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  transition="opacity 0.2s"
                >
                  {avatarMutation.isPending
                    ? <Spinner color="white" size="sm" />
                    : <Icon as={FiCamera} color="white" boxSize={5} />
                  }
                </Box>
              </Box>
              {user?.avatar_url && (
                <Box
                  as="button"
                  position="absolute"
                  bottom={0}
                  right={0}
                  w={6}
                  h={6}
                  borderRadius="full"
                  bg="red.500"
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  onClick={() => deleteAvatarMutation.mutate()}
                  title={t('profile.deleteAvatar')}
                  _hover={{ bg: 'red.600' }}
                  transition="background 0.15s"
                >
                  {deleteAvatarMutation.isPending
                    ? <Spinner color="white" size="xs" />
                    : <Icon as={FiTrash2} color="white" boxSize={3} />
                  }
                </Box>
              )}
            </Box>
            <input
              ref={avatarInputRef}
              type="file"
              accept="image/jpeg,image/png,image/webp"
              style={{ display: 'none' }}
              onChange={e => {
                const file = e.target.files?.[0]
                if (file) avatarMutation.mutate(file)
                e.target.value = ''
              }}
            />
            <VStack spacing={1.5} align="center">
              <Heading size="sm" letterSpacing="-0.2px">{fullName || user?.name}</Heading>
              <HStack spacing={1.5}>
                <Icon as={FiMail} color="faintText" boxSize={3} />
                <Text fontSize="xs" color="faintText">{user?.email}</Text>
              </HStack>
            </VStack>
          </VStack>

          <Divider my={5} borderColor="subtleBorder" />

          <VStack align="stretch" spacing={2.5}>
            {user?.city && (
              <HStack spacing={2.5}>
                <Icon as={FiMapPin} color="brand.400" boxSize={3.5} flexShrink={0} />
                <Text fontSize="sm" color="dimText">{user.city}</Text>
              </HStack>
            )}
            {user?.birth_date && (
              <HStack spacing={2.5}>
                <Icon as={FiCalendar} color="brand.400" boxSize={3.5} flexShrink={0} />
                <Text fontSize="sm" color="dimText">{formatDate(user.birth_date)}</Text>
              </HStack>
            )}
            {user?.created_at && (
              <HStack spacing={2.5}>
                <Icon as={FiUser} color="brand.400" boxSize={3.5} flexShrink={0} />
                <Text fontSize="xs" color="faintText">
                  {t('profile.memberSince')} {formatDate(user.created_at)}
                </Text>
              </HStack>
            )}
          </VStack>
        </Box>
      </Box>

      {/* Right panel */}
      <Box
        bg="cardBg"
        borderRadius="xl"
        border="1px solid"
        borderColor="cardBorder"
        boxShadow="0 1px 3px rgba(15,23,42,0.04), 0 4px 16px rgba(15,23,42,0.05)"
        overflow="hidden"
      >
        <Tabs colorScheme="brand" isLazy>
          <TabList
            px={2}
            pt={1}
            borderBottom="1px solid"
            borderColor="subtleBorder"
          >
            <Tab
              px={4} py={3} fontSize="sm" fontWeight="500" color="faintText"
              _selected={{ color: 'brand.600', fontWeight: '700', borderBottomColor: 'brand.500', borderBottomWidth: '2px' }}
              _hover={{ color: 'dimText' }}
            >
              {t('profile.editTab')}
            </Tab>
            <Tab
              px={4} py={3} fontSize="sm" fontWeight="500" color="faintText"
              _selected={{ color: 'brand.600', fontWeight: '700', borderBottomColor: 'brand.500', borderBottomWidth: '2px' }}
              _hover={{ color: 'dimText' }}
            >
              {t('profile.statsTab')}
            </Tab>
          </TabList>

          <TabPanels>
            {/* Edit tab */}
            <TabPanel p={{ base: 4, sm: 6 }}>
              <VStack spacing={5} align="stretch" maxW="480px">
                <FormControl>
                  <FormLabel fontSize="sm" color="dimText" mb={1.5}>{t('auth.email')}</FormLabel>
                  <Input value={user?.email ?? ''} isReadOnly opacity={0.6} cursor="default" _focus={{ boxShadow: 'none', borderColor: 'defaultBorder' }} />
                </FormControl>

                <Grid templateColumns={{ base: '1fr', sm: '1fr 1fr' }} gap={4}>
                  <FormControl>
                    <FormLabel fontSize="sm" color="dimText" mb={1.5}>{t('profile.firstName')}</FormLabel>
                    <Input
                      value={form.name}
                      onChange={e => setForm(f => ({ ...f, name: e.target.value }))}
                      placeholder={t('profile.firstName')}
                    />
                  </FormControl>
                  <FormControl>
                    <FormLabel fontSize="sm" color="dimText" mb={1.5}>{t('profile.lastName')}</FormLabel>
                    <Input
                      value={form.last_name}
                      onChange={e => setForm(f => ({ ...f, last_name: e.target.value }))}
                      placeholder={t('profile.lastName')}
                    />
                  </FormControl>
                </Grid>

                <FormControl>
                  <FormLabel fontSize="sm" color="dimText" mb={1.5}>{t('profile.city')}</FormLabel>
                  <Input
                    value={form.city}
                    onChange={e => setForm(f => ({ ...f, city: e.target.value }))}
                    placeholder={t('profile.city')}
                  />
                </FormControl>

                <FormControl>
                  <FormLabel fontSize="sm" color="dimText" mb={1.5}>{t('profile.birthDate')}</FormLabel>
                  <Input
                    type="date"
                    value={form.birth_date}
                    onChange={e => setForm(f => ({ ...f, birth_date: e.target.value }))}
                  />
                </FormControl>

                <Button
                  colorScheme="blue"
                  alignSelf="flex-start"
                  onClick={handleSave}
                  isLoading={updateMutation.isPending}
                  isDisabled={!form.name.trim()}
                >
                  {t('profile.save')}
                </Button>
              </VStack>
            </TabPanel>

            {/* Stats tab */}
            <TabPanel p={{ base: 4, sm: 6 }}>
              <Grid templateColumns={{ base: '1fr 1fr', md: 'repeat(4, 1fr)' }} gap={4}>
                <StatCard value={stats?.organized ?? 0} label={t('profile.statsOrganized')} />
                <StatCard value={stats?.participated ?? 0} label={t('profile.statsParticipated')} />
                <StatCard value={stats?.completed ?? 0} label={t('profile.statsCompleted')} />
                <StatCard value={stats?.upcoming ?? 0} label={t('profile.statsUpcoming')} />
              </Grid>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>
    </Grid>
  )
}

function StatCard({ value, label }: { value: number; label: string }) {
  return (
    <Box
      bg="subtleBg"
      borderRadius="xl"
      border="1px solid"
      borderColor="subtleBorder"
      p={5}
      textAlign="center"
    >
      <Text fontSize="3xl" fontWeight="700" color="brand.500" lineHeight={1}>
        {value}
      </Text>
      <Text fontSize="xs" color="faintText" mt={2} fontWeight="500">
        {label}
      </Text>
    </Box>
  )
}
