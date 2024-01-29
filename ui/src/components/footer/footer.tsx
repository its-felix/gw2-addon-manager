import {Box, SpaceBetween } from '@cloudscape-design/components';
import React from 'react';
import { useMobile } from '../util/state/common';
import classes from './footer.module.scss';

export function Gw2AddonManagerFooter() {
  const isMobile = useMobile();

  return (
    <footer id="gw2am-custom-footer" className={classes['gw2am-footer']}>
      <SpaceBetween size={'xs'} direction={'vertical'}>
        <SpaceBetween size={isMobile ? 'xs' : 'm'} direction={isMobile ? 'vertical' : 'horizontal'}>
          <Box variant={'span'}>© 2024 Felix.9127</Box>
        </SpaceBetween>
        <Box variant={'small'}>This site is not affiliated with ArenaNet, Guild Wars 2, or any of their partners. All copyrights reserved to their respective owners.</Box>
        <Box variant={'small'}>© ArenaNet LLC. All rights reserved. NCSOFT, ArenaNet, Guild Wars, Guild Wars 2, GW2, Guild Wars 2: Heart of Thorns, Guild Wars 2: Path of Fire, Guild Wars 2: End of Dragons, and Guild Wars 2: Secrets of the Obscure and all associated logos, designs, and composite marks are trademarks or registered trademarks of NCSOFT Corporation.</Box>
      </SpaceBetween>
    </footer>
  );
}
