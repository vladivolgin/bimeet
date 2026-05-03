import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Box,
  Button,
  Grid,
  HStack,
  Text,
  Badge,
  VStack,
  Icon,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  Avatar,
  Alert,
  AlertIcon,
  useDisclosure,
} from '@chakra-ui/react'
import { Link, useParams, useSearchParams } from 'react-router-dom'
import { useRef } from 'react'
import { useTranslation } from 'react-i18next'
import {
  FiCheckCircle,
  FiClock,
  FiUsers,
  FiPocket,
  FiBarChart2,
  FiShoppingBag,
  FiTruck,
  FiLink,
} from 'react-icons/fi'
import { useEvent } from '@/entities/event/queries'
import { useAuth } from '@/features/auth/model/AuthContext'
import { useCompleteEvent } from '@/features/event-manage/model/hooks'
import ParticipantsTab from '@/features/participants/ui/ParticipantsTab'
import CollectionsTab from '@/features/collections/ui/CollectionsTab'
import PollsTab from '@/features/polls/ui/PollsTab'
import ItemsTab from '@/features/items/ui/ItemsTab'
import CarpoolTab from '@/features/carpools/ui/CarpoolTab'
import LinksTab from '@/features/links/ui/LinksTab'
import EventInfoPanel from '@/widgets/event-detail/EventInfoPanel'
import ParticipantsSummaryPanel from '@/widgets/event-detail/ParticipantsSummaryPanel'
import { formatDateTime } from '@/shared/lib/formatDate'

const ORDINARY_TAB_KEYS = [
  { tKey: 'tabs.participants', icon: FiUsers,       key: 'participants' },
  { tKey: 'tabs.collections',  icon: FiPocket,      key: 'collections' },
  { tKey: 'tabs.polls',        icon: FiBarChart2,   key: 'polls' },
  { tKey: 'tabs.items',        icon: FiShoppingBag, key: 'items' },
  { tKey: 'tabs.carpools',     icon: FiTruck,       key: 'carpools' },
]

const BUSINESS_TAB_KEYS = [
  { tKey: 'tabs.participants', icon: FiUsers,     key: 'participants' },
  { tKey: 'tabs.polls',        icon: FiBarChart2, key: 'polls' },
  { tKey: 'tabs.links',        icon: FiLink,      key: 'links' },
]

export default function EventDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [searchParams] = useSearchParams()
  const event = useEvent(id!)
  const { user } = useAuth()
  const { t } = useTranslation()
  const { isOpen: isCompleteOpen, onOpen: onCompleteOpen, onClose: onCompleteClose } = useDisclosure()
  const cancelCompleteRef = useRef<HTMLButtonElement>(null)

  const completeEvent = useCompleteEvent(id ?? '')

  if (!event) {
    return (
      <Box textAlign="center" py={20}>
        <Text color="dimText" mb={4}>{t('events.notFound')}</Text>
        <Button as={Link} to="/events" colorScheme="blue">{t('common.backToList')}</Button>
      </Box>
    )
  }

  const isOrganizer = event.organizer.id === user?.id

  const baseTabKeys = event.category === 'business' ? BUSINESS_TAB_KEYS : ORDINARY_TAB_KEYS
  const allTabKeys = event.change_log.length > 0
    ? [...baseTabKeys, { tKey: 'tabs.history', icon: FiClock, key: 'changelog' }]
    : baseTabKeys
  const allTabs = allTabKeys.map(tab => ({ ...tab, label: t(tab.tKey) }))

  const tabParam = searchParams.get('tab')
  const defaultTabIndex = tabParam
    ? Math.max(0, allTabKeys.findIndex(t => t.key === tabParam))
    : 0

  return (
    <>
    <Grid
      templateColumns={{ base: '1fr', lg: '320px 1fr' }}
      gap={5}
      alignItems="flex-start"
    >
      {/* Left panel — sticky on desktop */}
      <Box position={{ lg: 'sticky' }} top={{ lg: '64px' }} zIndex={1}>
        <EventInfoPanel event={event} isOrganizer={isOrganizer} />
        <Box mt={4}>
          <ParticipantsSummaryPanel event={event} />
        </Box>
      </Box>

      {/* Right panel — tabs */}
      <Box
        bg="cardBg"
        borderRadius="xl"
        border="1px solid"
        borderColor="cardBorder"
        boxShadow="0 1px 3px rgba(15,23,42,0.04), 0 4px 16px rgba(15,23,42,0.05)"
        overflow="hidden"
      >
        {event.status === 'completed' && (
          <Alert status="warning" fontSize="sm" borderRadius={0} py={2.5}>
            <AlertIcon boxSize={4} />
            {t('events.completedAlert')}
          </Alert>
        )}

        {/* Mobile: complete button above tabs */}
        {isOrganizer && event.status === 'active' && (
          <Box px={4} pt={3} display={{ base: 'block', lg: 'none' }}>
            <Button
              w="full"
              size="sm"
              variant="outline"
              leftIcon={<FiCheckCircle size={14} />}
              borderColor="defaultBorder"
              color="dimText"
              fontWeight="500"
              _hover={{ borderColor: 'green.400', color: 'green.500', bg: 'transparent' }}
              _dark={{ _hover: { borderColor: 'green.400', color: 'green.400' } }}
              onClick={onCompleteOpen}
            >
              {t('eventInfo.menuComplete')}
            </Button>
          </Box>
        )}

        <Tabs colorScheme="brand" isLazy defaultIndex={defaultTabIndex}>
          <Box position="relative">
          <TabList
            px={2}
            pt={1}
            borderBottom="1px solid"
            borderColor="subtleBorder"
            overflowX={{ base: 'auto', lg: 'visible' }}
            css={{ scrollbarWidth: 'none', '&::-webkit-scrollbar': { display: 'none' } }}
            pr={{ lg: isOrganizer && event.status === 'active' ? '160px' : 2 }}
          >
            {allTabs.map(tab => (
              <Tab
                key={tab.key}
                px={4}
                py={3}
                fontSize="sm"
                fontWeight="500"
                color="faintText"
                whiteSpace="nowrap"
                _selected={{ color: 'brand.600', fontWeight: '700', borderBottomColor: 'brand.500', borderBottomWidth: '2px' }}
                _hover={{ color: 'dimText' }}
              >
                <HStack spacing={1.5}>
                  <Icon as={tab.icon} boxSize={3.5} />
                  <Text>{tab.label}</Text>
                </HStack>
              </Tab>
            ))}
          </TabList>

          {/* Desktop: complete button pinned to the right of the tab bar */}
          {isOrganizer && event.status === 'active' && (
            <Box
              position="absolute"
              right={3}
              top={0}
              bottom={0}
              display={{ base: 'none', lg: 'flex' }}
              alignItems="center"
            >
              <Button
                size="xs"
                variant="outline"
                leftIcon={<FiCheckCircle size={12} />}
                borderColor="defaultBorder"
                color="dimText"
                fontWeight="500"
                _hover={{ borderColor: 'green.400', color: 'green.500', bg: 'transparent' }}
                _dark={{ _hover: { borderColor: 'green.400', color: 'green.400' } }}
                onClick={onCompleteOpen}
              >
                {t('eventInfo.menuComplete')}
              </Button>
            </Box>
          )}
          </Box>


          <TabPanels>
            {allTabs.map(tab => (
              <TabPanel key={tab.key} p={{ base: 4, sm: 6 }}>
                {tab.key === 'participants' && <ParticipantsTab event={event} />}
                {tab.key === 'collections' && <CollectionsTab event={event} />}
                {tab.key === 'polls' && <PollsTab event={event} />}
                {tab.key === 'items' && <ItemsTab event={event} />}
                {tab.key === 'carpools' && <CarpoolTab event={event} />}
                {tab.key === 'links' && <LinksTab event={event} />}
                {tab.key === 'changelog' && (
                  <VStack align="stretch" spacing={3}>
                    <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.06em">
                      {t('events.changeHistory')}
                    </Text>
                    {[...event.change_log].reverse().map(log => (
                      <Box
                        key={log.id}
                        p={3.5}
                        borderRadius="xl"
                        border="1px solid"
                        borderColor="subtleBorder"
                        bg="subtleBg"
                      >
                        <HStack spacing={2.5} mb={2} flexWrap="wrap">
                          <Avatar size="xs" name={log.changed_by.name} bg="brand.400" flexShrink={0} />
                          <Text fontSize="sm" fontWeight="500">{log.changed_by.name}</Text>
                          <Text fontSize="xs" color="faintText">{t('events.changed', { field: log.field_name })}</Text>
                          <Text fontSize="xs" color="faintText" ml="auto">{formatDateTime(log.changed_at)}</Text>
                        </HStack>
                        <HStack spacing={2} flexWrap="wrap">
                          <Badge colorScheme="red" variant="subtle" borderRadius="md">{log.old_value}</Badge>
                          <Text color="faintText" fontSize="sm">→</Text>
                          <Badge colorScheme="green" variant="subtle" borderRadius="md">{log.new_value}</Badge>
                        </HStack>
                      </Box>
                    ))}
                  </VStack>
                )}
              </TabPanel>
            ))}
          </TabPanels>
        </Tabs>
      </Box>
    </Grid>

    <AlertDialog isOpen={isCompleteOpen} leastDestructiveRef={cancelCompleteRef} onClose={onCompleteClose}>
      <AlertDialogOverlay>
        <AlertDialogContent>
          <AlertDialogHeader fontSize="lg" fontWeight="bold">
            {t('eventInfo.completeDialogTitle')}
          </AlertDialogHeader>
          <AlertDialogBody>
            {t('eventInfo.completeDialogBody')}
          </AlertDialogBody>
          <AlertDialogFooter>
            <Button ref={cancelCompleteRef} onClick={onCompleteClose} variant="ghost">
              {t('common.cancel')}
            </Button>
            <Button
              colorScheme="blue"
              ml={3}
              onClick={() => completeEvent.mutate(undefined, { onSuccess: onCompleteClose })}
              isLoading={completeEvent.isPending}
            >
              {t('eventInfo.completeDialogConfirm')}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
    </>
  )
}
