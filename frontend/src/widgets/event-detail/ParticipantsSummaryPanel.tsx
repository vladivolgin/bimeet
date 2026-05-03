import {
  Box,
  HStack,
  Text,
  VStack,
  Avatar,
  AvatarGroup,
  Badge,
} from '@chakra-ui/react'
import { useTranslation } from 'react-i18next'
import type { Event } from '@/entities/event/model/types'

interface ParticipantsSummaryPanelProps {
  event: Event
}

export default function ParticipantsSummaryPanel({ event }: ParticipantsSummaryPanelProps) {
  const { t } = useTranslation()
  const confirmed = event.participants.filter(p => p.status === 'confirmed')
  const invited = event.participants.filter(p => p.status === 'invited')
  const declined = event.participants.filter(p => p.status === 'declined')

  if (event.participants.length === 0) return null

  return (
    <Box
      bg="cardBg"
      borderRadius="xl"
      border="1px solid"
      borderColor="cardBorder"
      boxShadow="0 1px 3px rgba(15,23,42,0.04), 0 4px 16px rgba(15,23,42,0.05)"
      p={4}
    >
      <Text
        fontSize="10px"
        fontWeight="600"
        color="faintText"
        textTransform="uppercase"
        letterSpacing="0.1em"
        mb={3}
      >
        {t('participantsSummary.title')}
      </Text>

      <VStack align="stretch" spacing={3}>
        <HStack spacing={2} flexWrap="wrap">
          <Badge colorScheme="green" variant="subtle" px={2.5} py={1}>
            {t('participantsSummary.going', { count: confirmed.length })}
          </Badge>
          <Badge colorScheme="yellow" variant="subtle" px={2.5} py={1}>
            {t('participantsSummary.waiting', { count: invited.length })}
          </Badge>
          {declined.length > 0 && (
            <Badge colorScheme="red" variant="subtle" px={2.5} py={1}>
              {t('participantsSummary.declined', { count: declined.length })}
            </Badge>
          )}
        </HStack>

        {confirmed.length > 0 && (
          <AvatarGroup size="sm" max={6} spacing="-8px">
            {confirmed.map(p => (
              <Avatar
                key={p.id}
                name={p.user.name}
                bg="brand.400"
                color="white"
                title={p.user.name}
              />
            ))}
          </AvatarGroup>
        )}
      </VStack>
    </Box>
  )
}
