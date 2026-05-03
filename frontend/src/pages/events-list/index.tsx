import {
  Box,
  Button,
  Collapse,
  Grid,
  Heading,
  HStack,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputLeftElement,
  ButtonGroup,
  Text,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  VStack,
  useDisclosure,
  useToast,
} from '@chakra-ui/react'
import { Link } from 'react-router-dom'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { FiCalendar, FiGrid, FiList, FiPlus, FiSearch } from 'react-icons/fi'
import { useEvents, usePublicEvents } from '@/entities/event/queries'
import { useAuth } from '@/features/auth/model/AuthContext'
import { useJoinPublicEvent, useConfirmAttendance } from '@/features/event-manage/model/hooks'
import EventCard from '@/widgets/event-card'
import EventListRow from '@/widgets/event-card/EventListRow'
import type { EventListItem, PublicEventListItem } from '@/entities/event/model/types'

type ViewMode = 'grid' | 'list'

function EmptyState({ onAction, label, isSearch }: { onAction?: string; label: string; isSearch?: boolean }) {
  const { t } = useTranslation()
  return (
    <Box textAlign="center" py={10}>
      <Icon as={isSearch ? FiSearch : FiCalendar} boxSize={8} color="faintText" mb={2} />
      <Text fontSize="sm" color="dimText">{label}</Text>
      {onAction && (
        <Button as={Link} to={onAction} colorScheme="blue" mt={4} size="sm">
          {t('events.emptyCreateButton')}
        </Button>
      )}
    </Box>
  )
}

function EventsGrid({
  events,
  myId,
  onAttend,
}: {
  events: (EventListItem | PublicEventListItem)[]
  myId?: string
  onAttend?: (event: EventListItem | PublicEventListItem) => void
}) {
  return (
    <Grid templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)', xl: 'repeat(3, 1fr)' }} gap={5}>
      {events.map(event => {
        const myParticipant = event.participants.find(p => p.user.id === myId)
        const myStatus = event.organizer.id !== myId
          ? (myParticipant?.status ?? event.my_status)
          : undefined
        const spotsLeft = event.max_guests != null ? event.max_guests - event.confirmed_count : null
        const canAttend = onAttend ? resolveCanAttend(event, myId, myStatus, spotsLeft) : undefined
        return (
          <EventCard
            key={event.id}
            event={event}
            myStatus={myStatus}
            canJoin={canAttend}
            onJoin={canAttend ? () => onAttend?.(event) : undefined}
          />
        )
      })}
    </Grid>
  )
}

function EventsList({
  events,
  myId,
  onAttend,
}: {
  events: (EventListItem | PublicEventListItem)[]
  myId?: string
  onAttend?: (event: EventListItem | PublicEventListItem) => void
}) {
  return (
    <VStack align="stretch" spacing={1.5}>
      {events.map(event => {
        const myParticipant = event.participants.find(p => p.user.id === myId)
        const myStatus = event.organizer.id !== myId
          ? (myParticipant?.status ?? event.my_status)
          : undefined
        const spotsLeft = event.max_guests != null ? event.max_guests - event.confirmed_count : null
        const canAttend = onAttend ? resolveCanAttend(event, myId, myStatus, spotsLeft) : undefined
        return (
          <EventListRow
            key={event.id}
            event={event}
            myStatus={myStatus}
            canJoin={canAttend}
            onJoin={canAttend ? () => onAttend?.(event) : undefined}
          />
        )
      })}
    </VStack>
  )
}

function resolveCanAttend(
  event: EventListItem | PublicEventListItem,
  myId: string | undefined,
  myStatus: string | undefined,
  spotsLeft: number | null,
): boolean | undefined {
  if (event.organizer.id === myId) return undefined
  if (myStatus === 'confirmed') return undefined
  if (myStatus === 'invited') return true
  if (event.is_public) {
    const isParticipant = (event as PublicEventListItem).is_participant
    if (isParticipant) return undefined
    return spotsLeft == null || spotsLeft > 0
  }
  return undefined
}

export default function EventsListPage() {
  const events = useEvents()
  const publicEvents = usePublicEvents()
  const { user } = useAuth()
  const { t } = useTranslation()
  const toast = useToast()
  const [search, setSearch] = useState('')
  const [viewMode, setViewMode] = useState<ViewMode>('grid')
  const { isOpen: isSearchOpen, onToggle: onSearchToggle } = useDisclosure()
  const joinPublicEvent = useJoinPublicEvent()
  const confirmAttendance = useConfirmAttendance()

  const allEvents = useMemo(() => {
    const map = new Map<string, EventListItem | PublicEventListItem>()
    events.forEach(e => map.set(e.id, e))
    publicEvents.forEach(e => { if (!map.has(e.id)) map.set(e.id, e) })
    return [...map.values()]
      .filter(e => e.status !== 'completed')
      .sort((a, b) => new Date(a.date_start).getTime() - new Date(b.date_start).getTime())
  }, [events, publicEvents])

  const attendingEvents = events.filter(e => e.organizer.id !== user?.id)
  const organizingEvents = events.filter(e => e.organizer.id === user?.id)

  function handleAttend(event: EventListItem | PublicEventListItem) {
    const myParticipant = event.participants.find(p => p.user.id === user?.id)
    const myStatus = myParticipant?.status ?? event.my_status

    if (myStatus === 'invited' && user?.id) {
      confirmAttendance.mutate(
        { eventId: event.id, userId: user.id },
        {
          onSuccess: () => toast({ title: t('eventCard.statusConfirmed'), status: 'success', duration: 2000 }),
          onError: (err) => toast({ title: t('common.error'), description: err.message, status: 'error', duration: 3000 }),
        },
      )
    } else {
      joinPublicEvent.mutate(event.id, {
        onSuccess: () => toast({ title: t('eventInfo.joinSuccess'), status: 'success', duration: 2000 }),
        onError: (err) => toast({ title: t('common.error'), description: err.message, status: 'error', duration: 3000 }),
      })
    }
  }

  function filterList(list: (EventListItem | PublicEventListItem)[]) {
    if (!search.trim()) return list
    const q = search.toLowerCase()
    return list.filter(e =>
      e.title.toLowerCase().includes(q) ||
      e.location.toLowerCase().includes(q)
    )
  }

  const filteredAll = filterList(allEvents)
  const filteredAttending = filterList(attendingEvents)
  const filteredOrganizing = filterList(organizingEvents)
  const hasSearch = search.trim().length > 0

  const tabStyle = {
    pb: 3, px: 4, fontSize: 'sm', fontWeight: '500', color: 'faintText',
    _selected: { color: 'brand.600', fontWeight: '700', borderBottomColor: 'brand.500', borderBottomWidth: '2px' },
  }

  function TabCount({ count }: { count: number }) {
    return (
      <Box as="span" ml={2} px={2} py={0.5} bg="brand.50" color="brand.600" borderRadius="full" fontSize="xs" fontWeight="700">
        {count}
      </Box>
    )
  }

  return (
    <Box>
      {/* Header */}
      <HStack justify="space-between" mb={4} align="center" gap={3}>
        <Heading size="lg" flexShrink={0}>{t('events.title')}</Heading>

        <HStack spacing={2} flexWrap="wrap" justify="flex-end">
          {/* Search — desktop only */}
          <InputGroup size="sm" w="220px" display={{ base: 'none', lg: 'flex' }}>
            <InputLeftElement pointerEvents="none">
              <Icon as={FiSearch} color="faintText" boxSize={3.5} />
            </InputLeftElement>
            <Input
              placeholder={t('events.searchPlaceholder')}
              value={search}
              onChange={e => setSearch(e.target.value)}
              borderRadius="10px"
              pl={8}
            />
          </InputGroup>

          {/* Search toggle — mobile */}
          <IconButton
            aria-label={t('events.search')}
            icon={<FiSearch size={15} />}
            variant="outline"
            size="sm"
            borderRadius="10px"
            onClick={onSearchToggle}
            display={{ base: 'flex', lg: 'none' }}
            color={isSearchOpen ? 'brand.600' : 'dimText'}
            borderColor={isSearchOpen ? 'brand.300' : 'defaultBorder'}
          />

          {/* View toggle */}
          <ButtonGroup size="sm" isAttached variant="outline">
            <IconButton
              aria-label={t('events.viewGrid')}
              icon={<FiGrid size={14} />}
              onClick={() => setViewMode('grid')}
              borderRadius="10px 0 0 10px"
              color={viewMode === 'grid' ? 'brand.600' : 'dimText'}
              borderColor={viewMode === 'grid' ? 'brand.300' : 'defaultBorder'}
              bg={viewMode === 'grid' ? 'navActiveBg' : undefined}
            />
            <IconButton
              aria-label={t('events.viewList')}
              icon={<FiList size={14} />}
              onClick={() => setViewMode('list')}
              borderRadius="0 10px 10px 0"
              color={viewMode === 'list' ? 'brand.600' : 'dimText'}
              borderColor={viewMode === 'list' ? 'brand.300' : 'defaultBorder'}
              bg={viewMode === 'list' ? 'navActiveBg' : undefined}
            />
          </ButtonGroup>

          <Button as={Link} to="/events/new" colorScheme="blue" leftIcon={<FiPlus />} size="sm">
            {t('events.createEvent')}
          </Button>
        </HStack>
      </HStack>

      {/* Mobile search bar */}
      <Collapse in={isSearchOpen} animateOpacity>
        <Box mb={4} display={{ base: 'block', lg: 'none' }}>
          <InputGroup size="sm">
            <InputLeftElement pointerEvents="none">
              <Icon as={FiSearch} color="faintText" boxSize={3.5} />
            </InputLeftElement>
            <Input
              placeholder={t('events.searchPlaceholder')}
              value={search}
              onChange={e => setSearch(e.target.value)}
              borderRadius="10px"
              pl={8}
              autoFocus
            />
          </InputGroup>
        </Box>
      </Collapse>

      <Tabs colorScheme="brand" variant="line">
        <TabList mb={6} borderColor="subtleBorder">
          <Tab {...tabStyle}>
            {t('events.tabAll')}
            <TabCount count={allEvents.length} />
          </Tab>
          <Tab {...tabStyle}>
            {t('events.tabAttending')}
            <TabCount count={attendingEvents.length} />
          </Tab>
          <Tab {...tabStyle}>
            {t('events.tabOrganizing')}
            <TabCount count={organizingEvents.length} />
          </Tab>
        </TabList>

        <TabPanels>
          {/* Все встречи */}
          <TabPanel px={0} pt={0}>
            {filteredAll.length === 0 ? (
              hasSearch
                ? <EmptyState label={t('events.emptySearchResult', { query: search })} isSearch />
                : <EmptyState onAction="/events/new" label={t('events.emptyNoOrganized')} />
            ) : viewMode === 'grid' ? (
              <EventsGrid events={filteredAll} myId={user?.id} onAttend={handleAttend} />
            ) : (
              <EventsList events={filteredAll} myId={user?.id} onAttend={handleAttend} />
            )}
          </TabPanel>

          {/* Участвую */}
          <TabPanel px={0} pt={0}>
            {filteredAttending.length === 0 ? (
              hasSearch
                ? <EmptyState label={t('events.emptySearchResult', { query: search })} isSearch />
                : <EmptyState label={t('events.emptyNoInvited')} />
            ) : viewMode === 'grid' ? (
              <EventsGrid events={filteredAttending} myId={user?.id} />
            ) : (
              <EventsList events={filteredAttending} myId={user?.id} />
            )}
          </TabPanel>

          {/* Организую */}
          <TabPanel px={0} pt={0}>
            {filteredOrganizing.length === 0 ? (
              hasSearch
                ? <EmptyState label={t('events.emptySearchResult', { query: search })} isSearch />
                : <EmptyState onAction="/events/new" label={t('events.emptyNoOrganized')} />
            ) : viewMode === 'grid' ? (
              <EventsGrid events={filteredOrganizing} myId={user?.id} />
            ) : (
              <EventsList events={filteredOrganizing} myId={user?.id} />
            )}
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Box>
  )
}
