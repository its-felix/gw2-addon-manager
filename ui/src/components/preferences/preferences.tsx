import {
  Box,
  Button,
  ColumnLayout,
  Header,
  Modal,
  ModalProps,
  SpaceBetween,
  Tiles,
} from '@cloudscape-design/components';
import React, { useEffect, useMemo, useState } from 'react';
import {
  ColorScheme, DateFormat, Preferences, UIDensity,
} from '../../lib/preferences.model';
import { ISO8601DateFormatter, SystemDateFormatter } from '../util/state/use-dateformat';
import { usePreferences } from '../util/state/use-preferences';

export function PreferencesModal(props: ModalProps) {
  const [preferences, setPreferences] = usePreferences();
  const [tempPreferences, setTempPreferences] = useState<Preferences>(preferences);
  const date = useMemo(() => new Date(), []);

  useEffect(() => {
    setTempPreferences(preferences);
  }, [preferences]);

  const { onDismiss } = props;
  function onCancelClick(e: CustomEvent) {
    setTempPreferences(preferences);

    if (onDismiss) {
      onDismiss(new CustomEvent(e.type, { detail: { reason: 'cancel' } }));
    }
  }

  function onSaveClick(e: CustomEvent) {
    setPreferences(tempPreferences);

    if (onDismiss) {
      onDismiss(new CustomEvent(e.type, { detail: { reason: 'save' } }));
    }
  }

  return (
    <Modal
      {...props}
      header={'Preferences'}
      size={'large'}
      footer={
        <Box float={'right'}>
          <SpaceBetween direction={'horizontal'} size={'xs'}>
            <Button variant={'link'} onClick={onCancelClick}>{'Cancel'}</Button>
            <Button variant={'primary'} onClick={onSaveClick}>{'Save'}</Button>
          </SpaceBetween>
        </Box>
      }
    >
      <ColumnLayout columns={1}>
        <div>
          <Header variant={'h3'}>Date and Time Format</Header>
          <Tiles
            value={tempPreferences.dateFormat}
            onChange={(e) => {
              setTempPreferences((prev) => ({ ...prev, dateFormat: e.detail.value as DateFormat }));
            }}
            items={[
              {
                label: 'System',
                description: SystemDateFormatter.formatDateTime(date),
                value: DateFormat.SYSTEM,
              },
              {
                label: 'ISO',
                description: ISO8601DateFormatter.formatDateTime(date),
                value: DateFormat.ISO_8601,
              },
            ]}
          />
        </div>
        <div>
          <Header variant={'h3'}>Color Scheme</Header>
          <Tiles
            value={tempPreferences.colorScheme}
            onChange={(e) => {
              setTempPreferences((prev) => ({ ...prev, colorScheme: e.detail.value as ColorScheme }));
            }}
            items={[
              {
                label: 'System',
                description: 'Use your system default color scheme',
                value: ColorScheme.SYSTEM,
              },
              {
                label: 'Light',
                description: 'Classic light theme',
                value: ColorScheme.LIGHT,
              },
              {
                label: 'Dark',
                description: 'Classic dark theme',
                value: ColorScheme.DARK,
              },
            ]}
          />
        </div>
        <div>
          <Header variant={'h3'}>Density</Header>
          <Tiles
            value={tempPreferences.uiDensity}
            onChange={(e) => {
              setTempPreferences((prev) => ({ ...prev, uiDensity: e.detail.value as UIDensity }));
            }}
            items={[
              {
                label: 'Comfortable',
                description: 'Standard spacing',
                value: UIDensity.COMFORTABLE,
              },
              {
                label: 'Compact',
                description: 'Reduced spacing',
                value: UIDensity.COMPACT,
              },
            ]}
          />
        </div>
      </ColumnLayout>
    </Modal>
  );
}
