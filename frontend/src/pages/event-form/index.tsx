import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  NumberInput,
  NumberInputField,
  Textarea,
  VStack,
  Heading,
  HStack,
  useToast,
  Text,
  SimpleGrid,
  RadioGroup,
  Radio,
  Stack,
  Divider,
} from '@chakra-ui/react'
import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { FiArrowLeft } from 'react-icons/fi'
import { useEvent } from '@/entities/event/queries'
import { useCreateEvent, useUpdateEvent } from '@/features/event-manage/model/hooks'
import type { EventCategory } from '@/entities/event/model/types'

interface FormValues {
  title: string
  description: string
  dateStart: string
  dateEnd: string
  location: string
  category: EventCategory
  dressCode: string
  isPublic: boolean
  maxGuests: string
}

function toDatetimeLocal(iso: string) {
  if (!iso) return ''
  return iso.slice(0, 16)
}

export default function EventFormPage() {
  const { id } = useParams<{ id?: string }>()
  const isEdit = Boolean(id)
  const navigate = useNavigate()
  const toast = useToast()
  const { t } = useTranslation()

  const event = useEvent(id ?? '')
  const createEvent = useCreateEvent()
  const updateEvent = useUpdateEvent(id ?? '')

  const [form, setForm] = useState<FormValues>({
    title: '',
    description: '',
    dateStart: '',
    dateEnd: '',
    location: '',
    category: 'ordinary',
    dressCode: '',
    isPublic: false,
    maxGuests: '',
  })

  useEffect(() => {
    if (event) {
      setForm({
        title: event.title,
        description: event.description,
        dateStart: toDatetimeLocal(event.date_start),
        dateEnd: toDatetimeLocal(event.date_end),
        location: event.location,
        category: event.category ?? 'ordinary',
        dressCode: event.dress_code ?? '',
        isPublic: event.is_public ?? false,
        maxGuests: event.max_guests != null ? String(event.max_guests) : '',
      })
    }
  }, [event])

  function handleChange(field: keyof FormValues, value: string) {
    setForm(f => ({ ...f, [field]: value }))
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    if (!form.title || !form.dateStart || !form.dateEnd || !form.location) {
      toast({ title: t('eventForm.requiredFieldsWarning'), status: 'warning', duration: 3000 })
      return
    }

    const maxGuests = form.isPublic && form.maxGuests.trim() ? parseInt(form.maxGuests) : undefined
    const payload = {
      title: form.title,
      description: form.description,
      date_start: new Date(form.dateStart).toISOString(),
      date_end: new Date(form.dateEnd).toISOString(),
      location: form.location,
      category: form.category,
      is_public: form.isPublic,
      ...(form.dressCode.trim() ? { dress_code: form.dressCode.trim() } : {}),
      ...(maxGuests != null ? { max_guests: maxGuests } : {}),
    }

    if (isEdit && id) {
      updateEvent.mutate(payload, {
        onSuccess: () => {
          toast({ title: t('eventForm.successUpdate'), status: 'success', duration: 2000 })
          navigate(`/events/${id}`)
        },
        onError: (err) => {
          toast({ title: t('eventForm.errorUpdate'), description: err.message, status: 'error', duration: 3000 })
        },
      })
    } else {
      createEvent.mutate(payload, {
        onSuccess: (newEvent) => {
          toast({ title: t('eventForm.successCreate'), status: 'success', duration: 2000 })
          navigate(`/events/${newEvent.id}`)
        },
        onError: (err) => {
          toast({ title: t('eventForm.errorCreate'), description: err.message, status: 'error', duration: 3000 })
        },
      })
    }
  }

  const isLoading = createEvent.isPending || updateEvent.isPending

  return (
    <Box maxW="860px" mx="auto">
      <Button
        as={Link}
        to={isEdit && id ? `/events/${id}` : '/events'}
        leftIcon={<FiArrowLeft />}
        variant="ghost"
        size="sm"
        color="dimText"
        mb={6}
        _hover={{ color: 'brand.600', bg: 'brand.50' }}
      >
        {t('common.back')}
      </Button>

      <Box mb={6}>
        <Heading size="lg">{isEdit ? t('eventForm.titleEdit') : t('eventForm.titleCreate')}</Heading>
        <Text fontSize="sm" color="dimText" mt={1}>
          {isEdit ? t('eventForm.subtitleEdit') : t('eventForm.subtitleCreate')}
        </Text>
      </Box>

      <Box
        bg="cardBg"
        borderRadius="xl"
        border="1px solid"
        borderColor="cardBorder"
        boxShadow="0 1px 2px rgba(0,0,0,0.04), 0 4px 16px rgba(0,0,0,0.04)"
        overflow="hidden"
      >
        <Box h="4px" bgGradient="linear(135deg, brand.600, #7C3AED)" />

        <form onSubmit={handleSubmit}>
          <Box p={{ base: 4, sm: 6 }}>
            {/* Тип встречи — full width */}
            <Box mb={6}>
              <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.08em" mb={3}>
                {t('eventForm.sectionType')}
              </Text>
              <FormControl isRequired>
                <RadioGroup value={form.category} onChange={(val) => handleChange('category', val)}>
                  <Stack direction={{ base: 'column', sm: 'row' }} spacing={4}>
                    <Radio value="ordinary" colorScheme="blue">
                      <VStack align="flex-start" spacing={0}>
                        <Text fontSize="sm" fontWeight="500">{t('eventForm.typeOrdinary')}</Text>
                        <Text fontSize="xs" color="dimText">{t('eventForm.typeOrdinaryDesc')}</Text>
                      </VStack>
                    </Radio>
                    <Radio value="business" colorScheme="blue">
                      <VStack align="flex-start" spacing={0}>
                        <Text fontSize="sm" fontWeight="500">{t('eventForm.typeBusiness')}</Text>
                        <Text fontSize="xs" color="dimText">{t('eventForm.typeBusinessDesc')}</Text>
                      </VStack>
                    </Radio>
                  </Stack>
                </RadioGroup>
              </FormControl>
            </Box>

            <Divider borderColor="subtleBorder" my={6} />

            {/* Visibility */}
            <Box mb={6}>
              <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.08em" mb={3}>
                {t('eventForm.sectionVisibility')}
              </Text>
              <RadioGroup value={form.isPublic ? 'public' : 'personal'} onChange={(val) => setForm(f => ({ ...f, isPublic: val === 'public', maxGuests: val === 'personal' ? '' : f.maxGuests }))}>
                <Stack direction={{ base: 'column', sm: 'row' }} spacing={4}>
                  <Radio value="personal" colorScheme="blue">
                    <VStack align="flex-start" spacing={0}>
                      <Text fontSize="sm" fontWeight="500">{t('eventForm.typePersonal')}</Text>
                      <Text fontSize="xs" color="dimText">{t('eventForm.typePersonalDesc')}</Text>
                    </VStack>
                  </Radio>
                  <Radio value="public" colorScheme="blue">
                    <VStack align="flex-start" spacing={0}>
                      <Text fontSize="sm" fontWeight="500">{t('eventForm.typePublic')}</Text>
                      <Text fontSize="xs" color="dimText">{t('eventForm.typePublicDesc')}</Text>
                    </VStack>
                  </Radio>
                </Stack>
              </RadioGroup>
              {form.isPublic && (
                <FormControl mt={4} maxW="200px">
                  <FormLabel fontSize="sm" fontWeight="500" color="dimText">
                    {t('eventForm.maxGuests')}
                    <Text as="span" fontSize="xs" color="faintText" fontWeight="400" ml={2}>{t('common.optional')}</Text>
                  </FormLabel>
                  <NumberInput
                    min={1}
                    value={form.maxGuests}
                    onChange={(val) => setForm(f => ({ ...f, maxGuests: val }))}
                  >
                    <NumberInputField placeholder={t('eventForm.maxGuestsPlaceholder')} />
                  </NumberInput>
                </FormControl>
              )}
            </Box>

            <Divider borderColor="subtleBorder" mb={6} />

            {/* Two-column layout on md+ */}
            <SimpleGrid columns={{ base: 1, md: 2 }} spacing={6}>
              {/* Left column: main info */}
              <VStack spacing={4} align="stretch">
                <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.08em">
                  {t('eventForm.sectionMainInfo')}
                </Text>

                <FormControl isRequired>
                  <FormLabel fontSize="sm" fontWeight="500" color="dimText">{t('eventForm.fieldTitle')}</FormLabel>
                  <Input
                    value={form.title}
                    onChange={e => handleChange('title', e.target.value)}
                    placeholder={t('eventForm.fieldTitlePlaceholder')}
                  />
                </FormControl>

                <FormControl>
                  <FormLabel fontSize="sm" fontWeight="500" color="dimText">{t('eventForm.fieldDescription')}</FormLabel>
                  <Textarea
                    value={form.description}
                    onChange={e => handleChange('description', e.target.value)}
                    placeholder={t('eventForm.fieldDescriptionPlaceholder')}
                    rows={4}
                    resize="none"
                  />
                </FormControl>
              </VStack>

              {/* Right column: when & where */}
              <VStack spacing={4} align="stretch">
                <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.08em">
                  {t('eventForm.sectionWhenWhere')}
                </Text>

                <SimpleGrid columns={2} spacing={3}>
                  <FormControl isRequired>
                    <FormLabel fontSize="sm" fontWeight="500" color="dimText">{t('eventForm.fieldStart')}</FormLabel>
                    <Input
                      type="datetime-local"
                      value={form.dateStart}
                      onChange={e => handleChange('dateStart', e.target.value)}
                    />
                  </FormControl>
                  <FormControl isRequired>
                    <FormLabel fontSize="sm" fontWeight="500" color="dimText">{t('eventForm.fieldEnd')}</FormLabel>
                    <Input
                      type="datetime-local"
                      value={form.dateEnd}
                      onChange={e => handleChange('dateEnd', e.target.value)}
                    />
                  </FormControl>
                </SimpleGrid>

                <FormControl isRequired>
                  <FormLabel fontSize="sm" fontWeight="500" color="dimText">{t('eventForm.fieldLocation')}</FormLabel>
                  <Input
                    value={form.location}
                    onChange={e => handleChange('location', e.target.value)}
                    placeholder={form.category === 'business' ? t('eventForm.fieldLocationPlaceholderBusiness') : t('eventForm.fieldLocationPlaceholderOrdinary')}
                  />
                </FormControl>

                <FormControl>
                  <FormLabel fontSize="sm" fontWeight="500" color="dimText">
                    {t('eventForm.fieldDressCode')}
                    <Text as="span" fontSize="xs" color="faintText" fontWeight="400" ml={2}>{t('common.optional')}</Text>
                  </FormLabel>
                  <Input
                    value={form.dressCode}
                    onChange={e => handleChange('dressCode', e.target.value)}
                    placeholder={t('eventForm.fieldDressCodePlaceholder')}
                  />
                </FormControl>
              </VStack>
            </SimpleGrid>
          </Box>

          {/* Sticky submit bar */}
          <Box
            borderTop="1px solid"
            borderColor="subtleBorder"
            px={{ base: 4, sm: 6 }}
            py={4}
            bg="cardBg"
          >
            <HStack justify="space-between" align="center">
              {!isEdit && (
                <Text fontSize="xs" color="faintText">
                  {t('eventForm.footerHint')}
                </Text>
              )}
              <HStack spacing={2} ml="auto">
                <Button variant="ghost" onClick={() => navigate(-1)}>{t('common.cancel')}</Button>
                <Button type="submit" colorScheme="blue" px={8} isLoading={isLoading}>
                  {isEdit ? t('common.save') : t('common.create')}
                </Button>
              </HStack>
            </HStack>
          </Box>
        </form>
      </Box>
    </Box>
  )
}
