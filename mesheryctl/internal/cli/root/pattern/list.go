// Copyright 2023 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pattern

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbose bool
)

var linkDocPatternList = map[string]string{
	"link":    "![pattern-list-usage](/assets/img/mesheryctl/patternList.png)",
	"caption": "Usage of mesheryctl pattern list",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List patterns",
	Long:  `Display list of all available pattern files.`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// list all available patterns
mesheryctl pattern list
	`,
	Annotations: linkDocPatternList,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		var response models.PatternsAPIResponse
		client := &http.Client{}
		req, err := utils.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/pattern", nil)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return errors.Errorf("Unable to reach Meshery server at %s. Verify your environment's readiness for a Meshery deployment by running `mesheryctl system check`. ", mctlCfg.GetBaseMesheryURL())
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &response)

		if err != nil {
			return err
		}
		tokenObj, err := utils.ReadToken(utils.TokenFlag)
		if err != nil {
			return err
		}
		provider := tokenObj["meshery-provider"]
		var data [][]string

		if verbose {
			if provider == "None" {
				for _, v := range response.Patterns {
					PatternID := v.ID.String()
					PatterName := v.Name
					CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
					UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
					data = append(data, []string{PatternID, PatterName, CreatedAt, UpdatedAt})
				}
				utils.PrintToTableWithFooter([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", ""})
				return nil
			}

			for _, v := range response.Patterns {
				PatternID := utils.TruncateID(v.ID.String())
				var UserID string
				if v.UserID != nil {
					UserID = *v.UserID
				} else {
					UserID = "null"
				}
				PatterName := v.Name
				CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
				UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
				data = append(data, []string{PatternID, UserID, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTableWithFooter([]string{"PATTERN ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", "", ""})

			return nil
		}

		// Check if messhery provider is set
		if provider == "None" {
			for _, v := range response.Patterns {
				PatterName := strings.Trim(v.Name, filepath.Ext(v.Name))
				PatternID := utils.TruncateID(v.ID.String())
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternID, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTableWithFooter([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", ""})
			return nil
		}
		for _, v := range response.Patterns {
			PatternID := utils.TruncateID(v.ID.String())
			var UserID string
			if v.UserID != nil {
				UserID = *v.UserID
			} else {
				UserID = "null"
			}
			PatterName := v.Name
			CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
			UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
			data = append(data, []string{PatternID, UserID, PatterName, CreatedAt, UpdatedAt})
		}
		utils.PrintToTableWithFooter([]string{"PATTERN ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data, []string{"Total", fmt.Sprintf("%d", response.TotalCount), "", "", ""})

		return nil

	},
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display full length user and pattern file identifiers")
}
