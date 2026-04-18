/**
 * @Company: 云南奇讯科技有限公司
 * @Author: yxf
 * @Description:
 * @Date: 2024/1/25 16:22
 */

package eventBus

import (
	"gfast/internal/app/common/service"

	"github.com/asaskevich/EventBus"
)

func init() {
	service.RegisterEventBus(EventBus.New())
}
