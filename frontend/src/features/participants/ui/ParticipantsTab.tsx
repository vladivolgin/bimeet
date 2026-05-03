import {
  VStack,
  HStack,
  Stack,
  Avatar,
  Text,
  Button,
  Box,
  Badge,
  Input,
  FormControl,
  useDisclosure,
  Collapse,
  useToast,
  Icon,
} from '@chakra-ui/react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { FiPlus, FiLink } from 'react-icons/fi'
import type { Event } from '@/entities/event/model/types'
import { useInviteParticipant, useUpdateParticipantStatus } from '@/features/participants/model/hooks'
import { useAuth } from '@/features/auth/model/AuthContext'

// statusLabel is now handled via t() in the component

const statusScheme = {
  invited:   'yellow',
  confirmed: 'green',
  declined:  'red',
} as const

export default function ParticipantsTab({ event }: { event: Event }) {
  const { user } = useAuth()
  const toast = useToast()
  const { t } = useTranslation()
  const { isOpen, onToggle } = useDisclosure()
  const [email, setEmail] = useState('')

  const updateStatus = useUpdateParticipantStatus(event.id)
  const invite = useInviteParticipant(event.id)

  const isOrganizer = event.organizer.id === user?.id
  const myParticipant = event.participants.find(p => p.user.id === user?.id)

  const confirmed = event.participants.filter(p => p.status === 'confirmed').length
  const invited   = event.participants.filter(p => p.status === 'invited').length
  const declined  = event.participants.filter(p => p.status === 'declined').length

  function handleInvite() {
    if (!email.trim()) return
    invite.mutate(email.trim(), {
      onSuccess: () => {
        toast({ title: t('participants.inviteSent'), status: 'success', duration: 2000 })
        setEmail('')
        onToggle()
      },
      onError: (err) => {
        toast({ title: t('common.error'), description: err.message, status: 'error', duration: 3000 })
      },
    })
  }

  function handleCopyLink() {
    const link = `${window.location.origin}/invite/${event.invite_token}`
    navigator.clipboard.writeText(link).then(() => {
      toast({ title: t('participants.linkCopied'), status: 'success', duration: 2000 })
    })
  }

  const isCompleted = event.status === 'completed'

  return (
    <VStack align="stretch" spacing={4}>
      {myParticipant && !isOrganizer && !isCompleted && (
        <Box p={4} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder">
          <Text fontSize="sm" fontWeight="500" color="dimText" mb={3}>{t('participants.myStatusTitle')}</Text>
          <HStack spacing={2}>
            <Button
              flex={1}
              size="sm"
              colorScheme="green"
              variant={myParticipant.status === 'confirmed' ? 'solid' : 'outline'}
              onClick={() => user && updateStatus.mutate({ userId: user.id, status: 'confirmed' })}
              isLoading={updateStatus.isPending}
            >
              {t('participants.going')}
            </Button>
            <Button
              flex={1}
              size="sm"
              colorScheme="red"
              variant={myParticipant.status === 'declined' ? 'solid' : 'outline'}
              onClick={() => user && updateStatus.mutate({ userId: user.id, status: 'declined' })}
              isLoading={updateStatus.isPending}
            >
              {t('participants.notGoing')}
            </Button>
          </HStack>
        </Box>
      )}

      {isOrganizer && !isCompleted && (
        <Stack direction={{ base: 'column', sm: 'row' }} spacing={2}>
          <Button
            size="sm"
            leftIcon={<FiPlus />}
            colorScheme="blue"
            variant="outline"
            onClick={onToggle}
            flex={1}
          >
            {t('participants.inviteByEmail')}
          </Button>
          <Button
            size="sm"
            leftIcon={<Icon as={FiLink} />}
            variant="outline"
            onClick={handleCopyLink}
            flex={1}
          >
            {t('participants.copyLink')}
          </Button>
        </Stack>
      )}

      <Collapse in={isOpen} animateOpacity>
        <Box p={4} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder">
          <FormControl>
            <Text fontSize="xs" fontWeight="600" color="dimText" mb={2}>{t('participants.emailLabel')}</Text>
            <VStack spacing={2} align="stretch">
              <Input
                size="sm"
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                placeholder="user@example.com"
                onKeyDown={e => e.key === 'Enter' && handleInvite()}
                autoFocus
              />
              <HStack spacing={2} justify="flex-end">
                <Button size="sm" variant="ghost" onClick={onToggle}>{t('common.cancel')}</Button>
                <Button
                  size="sm"
                  colorScheme="blue"
                  onClick={handleInvite}
                  isLoading={invite.isPending}
                >
                  {t('participants.inviteButton')}
                </Button>
              </HStack>
            </VStack>
          </FormControl>
        </Box>
      </Collapse>

      <HStack spacing={2} flexWrap="wrap">
        <Badge colorScheme="green" variant="subtle" px={2.5} py={1}>{t('participants.badgeCount_going', { count: confirmed })}</Badge>
        <Badge colorScheme="yellow" variant="subtle" px={2.5} py={1}>{t('participants.badgeCount_invited', { count: invited })}</Badge>
        {declined > 0 && (
          <Badge colorScheme="red" variant="subtle" px={2.5} py={1}>{t('participants.badgeCount_declined', { count: declined })}</Badge>
        )}
      </HStack>

      <VStack align="stretch" spacing={2}>
        {event.participants.map(p => (
          <HStack key={p.id} justify="space-between" p={3} borderRadius="lg" bg="subtleBg">
            <HStack spacing={3} minW={0} flex={1}>
              <Avatar
                size="sm"
                name={p.user.name}
                bg="brand.100"
                color="brand.700"
                flexShrink={0}
              />
              <Box minW={0}>
                <HStack spacing={1.5}>
                  <Text fontWeight="500" fontSize="sm" noOfLines={1}>{p.user.name}</Text>
                  {p.user.id === event.organizer.id && (
                    <Badge colorScheme="purple" fontSize="xs" flexShrink={0}>{t('participants.badgeOrganizer')}</Badge>
                  )}
                </HStack>
                <Text fontSize="xs" color="faintText" noOfLines={1}>{p.user.email}</Text>
              </Box>
            </HStack>
            <Badge
              colorScheme={statusScheme[p.status]}
              variant="subtle"
              flexShrink={0}
            >
              {t(`participants.status${p.status.charAt(0).toUpperCase()}${p.status.slice(1)}`)}
            </Badge>
          </HStack>
        ))}

      </VStack>
    </VStack>
  )
}
